package superclaude

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/opencode-ai/opencode/internal/logging"
)

// Optimizer provides performance optimizations for SuperClaude
type Optimizer struct {
	// Response caching
	cache      sync.Map
	cacheSize  int
	cacheTTL   time.Duration
	
	// Request batching
	batchQueue chan *OptimizedRequest
	batchSize  int
	batchDelay time.Duration
	
	// Resource pooling
	workerPool *WorkerPool
	
	// Metrics
	metrics *Metrics
}

// OptimizedRequest wraps a request with optimization metadata
type OptimizedRequest struct {
	Command   string
	SessionID string
	Context   context.Context
	Response  chan *OptimizedResponse
	Timestamp time.Time
}

// OptimizedResponse contains the response and metrics
type OptimizedResponse struct {
	Result    interface{}
	Error     error
	CacheHit  bool
	BatchSize int
	Duration  time.Duration
}

// WorkerPool manages a pool of workers for parallel processing
type WorkerPool struct {
	workers   int
	taskQueue chan func()
	wg        sync.WaitGroup
}

// Metrics tracks performance metrics
type Metrics struct {
	mu            sync.RWMutex
	totalRequests int64
	cacheHits     int64
	avgDuration   time.Duration
	peakMemory    uint64
}

// NewOptimizer creates a new optimizer
func NewOptimizer() *Optimizer {
	opt := &Optimizer{
		cacheSize:  1000,
		cacheTTL:   15 * time.Minute,
		batchSize:  10,
		batchDelay: 100 * time.Millisecond,
		batchQueue: make(chan *OptimizedRequest, 100),
		metrics:    &Metrics{},
	}
	
	// Initialize worker pool based on CPU cores
	numWorkers := runtime.NumCPU() * 2
	opt.workerPool = NewWorkerPool(numWorkers)
	
	// Start batch processor
	go opt.processBatches()
	
	// Start cache cleaner
	go opt.cleanCache()
	
	// Start metrics collector
	go opt.collectMetrics()
	
	return opt
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	wp := &WorkerPool{
		workers:   workers,
		taskQueue: make(chan func(), workers*2),
	}
	
	// Start workers
	for i := 0; i < workers; i++ {
		go wp.worker()
	}
	
	return wp
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker() {
	for task := range wp.taskQueue {
		task()
	}
}

// Submit submits a task to the worker pool
func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1)
	wp.taskQueue <- func() {
		defer wp.wg.Done()
		task()
	}
}

// Wait waits for all tasks to complete
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// OptimizeCommand optimizes a SuperClaude command execution
func (opt *Optimizer) OptimizeCommand(ctx context.Context, sessionID, command string) (*OptimizedResponse, error) {
	start := time.Now()
	
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", sessionID, command)
	if cached, ok := opt.cache.Load(cacheKey); ok {
		if entry, ok := cached.(*CacheEntry); ok && !entry.IsExpired() {
			opt.recordCacheHit()
			return &OptimizedResponse{
				Result:   entry.Data,
				CacheHit: true,
				Duration: time.Since(start),
			}, nil
		}
	}
	
	// Create optimized request
	req := &OptimizedRequest{
		Command:   command,
		SessionID: sessionID,
		Context:   ctx,
		Response:  make(chan *OptimizedResponse, 1),
		Timestamp: time.Now(),
	}
	
	// Try to batch with other requests
	select {
	case opt.batchQueue <- req:
		// Added to batch queue
	case <-time.After(opt.batchDelay):
		// Process immediately if queue is full
		opt.processSingleRequest(req)
	}
	
	// Wait for response
	select {
	case resp := <-req.Response:
		// Cache successful responses
		if resp.Error == nil && !resp.CacheHit {
			opt.cache.Store(cacheKey, &CacheEntry{
				Data:      resp.Result,
				Timestamp: time.Now(),
			})
		}
		
		opt.recordRequest(time.Since(start))
		return resp, nil
		
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// processBatches processes requests in batches
func (opt *Optimizer) processBatches() {
	ticker := time.NewTicker(opt.batchDelay)
	defer ticker.Stop()
	
	batch := make([]*OptimizedRequest, 0, opt.batchSize)
	
	for {
		select {
		case req := <-opt.batchQueue:
			batch = append(batch, req)
			
			// Process batch if full
			if len(batch) >= opt.batchSize {
				opt.processBatch(batch)
				batch = batch[:0]
			}
			
		case <-ticker.C:
			// Process partial batch
			if len(batch) > 0 {
				opt.processBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// processBatch processes a batch of requests in parallel
func (opt *Optimizer) processBatch(batch []*OptimizedRequest) {
	logging.Debug("Processing batch", "size", len(batch))
	
	// Group by command type for better batching
	groups := make(map[string][]*OptimizedRequest)
	for _, req := range batch {
		cmdType := extractCommandType(req.Command)
		groups[cmdType] = append(groups[cmdType], req)
	}
	
	// Process each group in parallel
	var wg sync.WaitGroup
	for cmdType, reqs := range groups {
		wg.Add(1)
		go func(ct string, requests []*OptimizedRequest) {
			defer wg.Done()
			opt.processCommandGroup(ct, requests)
		}(cmdType, reqs)
	}
	
	wg.Wait()
}

// processCommandGroup processes a group of similar commands
func (opt *Optimizer) processCommandGroup(cmdType string, requests []*OptimizedRequest) {
	// Special optimization for certain command types
	switch cmdType {
	case "analyze":
		// Combine analyze requests for the same directory
		opt.combineAnalyzeRequests(requests)
	case "test":
		// Run tests in parallel
		opt.parallelTestRequests(requests)
	default:
		// Process individually
		for _, req := range requests {
			opt.processSingleRequest(req)
		}
	}
}

// processSingleRequest processes a single request
func (opt *Optimizer) processSingleRequest(req *OptimizedRequest) {
	// This would call the actual SuperClaude handler
	// For now, we'll simulate it
	result := fmt.Sprintf("Processed: %s", req.Command)
	
	req.Response <- &OptimizedResponse{
		Result:    result,
		Error:     nil,
		CacheHit:  false,
		BatchSize: 1,
		Duration:  time.Since(req.Timestamp),
	}
}

// combineAnalyzeRequests combines multiple analyze requests
func (opt *Optimizer) combineAnalyzeRequests(requests []*OptimizedRequest) {
	// Group by target directory
	targets := make(map[string][]*OptimizedRequest)
	for _, req := range requests {
		target := extractTarget(req.Command)
		targets[target] = append(targets[target], req)
	}
	
	// Analyze each target once and share results
	for target, reqs := range targets {
		result := fmt.Sprintf("Combined analysis of %s for %d requests", target, len(reqs))
		
		// Send result to all requests
		for _, req := range reqs {
			req.Response <- &OptimizedResponse{
				Result:    result,
				Error:     nil,
				CacheHit:  false,
				BatchSize: len(reqs),
				Duration:  time.Since(req.Timestamp),
			}
		}
	}
}

// parallelTestRequests runs test requests in parallel
func (opt *Optimizer) parallelTestRequests(requests []*OptimizedRequest) {
	var wg sync.WaitGroup
	
	for _, req := range requests {
		wg.Add(1)
		opt.workerPool.Submit(func() {
			defer wg.Done()
			opt.processSingleRequest(req)
		})
	}
	
	wg.Wait()
}

// CacheEntry represents a cached response
type CacheEntry struct {
	Data      interface{}
	Timestamp time.Time
}

// IsExpired checks if the cache entry is expired
func (ce *CacheEntry) IsExpired() bool {
	return time.Since(ce.Timestamp) > 15*time.Minute
}

// cleanCache periodically cleans expired cache entries
func (opt *Optimizer) cleanCache() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	
	for range ticker.C {
		count := 0
		opt.cache.Range(func(key, value interface{}) bool {
			if entry, ok := value.(*CacheEntry); ok && entry.IsExpired() {
				opt.cache.Delete(key)
				count++
			}
			return true
		})
		
		if count > 0 {
			logging.Debug("Cleaned cache entries", "count", count)
		}
	}
}

// collectMetrics periodically collects performance metrics
func (opt *Optimizer) collectMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		opt.metrics.mu.Lock()
		if m.Alloc > opt.metrics.peakMemory {
			opt.metrics.peakMemory = m.Alloc
		}
		opt.metrics.mu.Unlock()
		
		logging.Debug("Performance metrics",
			"total_requests", opt.metrics.totalRequests,
			"cache_hits", opt.metrics.cacheHits,
			"cache_hit_rate", opt.getCacheHitRate(),
			"avg_duration", opt.metrics.avgDuration,
			"memory_mb", m.Alloc/1024/1024,
			"goroutines", runtime.NumGoroutine(),
		)
	}
}

// Helper methods

func (opt *Optimizer) recordRequest(duration time.Duration) {
	opt.metrics.mu.Lock()
	defer opt.metrics.mu.Unlock()
	
	opt.metrics.totalRequests++
	
	// Update average duration
	if opt.metrics.avgDuration == 0 {
		opt.metrics.avgDuration = duration
	} else {
		opt.metrics.avgDuration = (opt.metrics.avgDuration + duration) / 2
	}
}

func (opt *Optimizer) recordCacheHit() {
	opt.metrics.mu.Lock()
	defer opt.metrics.mu.Unlock()
	opt.metrics.cacheHits++
}

func (opt *Optimizer) getCacheHitRate() float64 {
	opt.metrics.mu.RLock()
	defer opt.metrics.mu.RUnlock()
	
	if opt.metrics.totalRequests == 0 {
		return 0
	}
	
	return float64(opt.metrics.cacheHits) / float64(opt.metrics.totalRequests)
}

func extractCommandType(command string) string {
	parts := strings.Split(command, " ")
	if len(parts) > 0 && strings.HasPrefix(parts[0], "/user:") {
		return strings.TrimPrefix(parts[0], "/user:")
	}
	return "unknown"
}

func extractTarget(command string) string {
	parts := strings.Split(command, " ")
	if len(parts) > 1 {
		return parts[1]
	}
	return "."
}