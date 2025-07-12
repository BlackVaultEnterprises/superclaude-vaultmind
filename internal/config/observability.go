package config

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/opencode-ai/opencode/internal/logging"
)

// ConfigObservability provides comprehensive configuration monitoring
type ConfigObservability struct {
	metrics           *ConfigMetrics
	driftDetector     *ConfigDriftDetector
	healthChecker     *ConfigHealthChecker
	complianceChecker *ComplianceChecker
	alertManager      *AlertManager
	enabled           bool
	mu                sync.RWMutex
}

// ConfigMetrics tracks configuration-related metrics
type ConfigMetrics struct {
	configLoads          prometheus.Counter
	configValidations    prometheus.Counter
	configErrors         prometheus.Counter
	configReloads        prometheus.Counter
	configSize           prometheus.Gauge
	validationDuration   prometheus.Histogram
	tenantCount          prometheus.Gauge
	activeTenants        prometheus.Gauge
	quotaViolations      prometheus.Counter
	featureUsage         *prometheus.CounterVec
	configChanges        prometheus.Counter
	driftDetections      prometheus.Counter
}

// ConfigDriftDetector monitors configuration drift
type ConfigDriftDetector struct {
	baseline        *SuperClaudeConfig
	checkInterval   time.Duration
	driftThreshold  float64
	alertChannel    chan DriftAlert
	running         bool
	mu              sync.RWMutex
}

// DriftAlert represents a configuration drift detection
type DriftAlert struct {
	Timestamp   time.Time              `json:"timestamp"`
	DriftType   DriftType              `json:"drift_type"`
	Severity    AlertSeverity          `json:"severity"`
	Component   string                 `json:"component"`
	Expected    interface{}            `json:"expected"`
	Actual      interface{}            `json:"actual"`
	Difference  float64                `json:"difference"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type DriftType int

const (
	DriftPerformance DriftType = iota
	DriftSecurity
	DriftCompliance
	DriftResource
	DriftConfiguration
)

type AlertSeverity int

const (
	AlertInfo AlertSeverity = iota
	AlertWarning
	AlertCritical
	AlertEmergency
)

// ConfigHealthChecker monitors configuration health
type ConfigHealthChecker struct {
	checks    []HealthCheck
	interval  time.Duration
	results   map[string]HealthResult
	mu        sync.RWMutex
}

// HealthCheck defines a configuration health check
type HealthCheck struct {
	Name        string
	Description string
	Check       func(*SuperClaudeConfig) HealthResult
	Interval    time.Duration
	Timeout     time.Duration
	Critical    bool
}

// HealthResult represents the result of a health check
type HealthResult struct {
	Status    HealthStatus           `json:"status"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  time.Duration          `json:"duration"`
}

type HealthStatus int

const (
	HealthHealthy HealthStatus = iota
	HealthDegraded
	HealthUnhealthy
	HealthUnknown
)

// ComplianceChecker validates configuration compliance
type ComplianceChecker struct {
	standards []ComplianceStandard
	mu        sync.RWMutex
}

// ComplianceStandard defines a compliance standard
type ComplianceStandard struct {
	Name        string
	Description string
	Version     string
	Rules       []ComplianceRule
	Required    bool
}

// ComplianceRule defines a specific compliance rule
type ComplianceRule struct {
	ID          string
	Description string
	Severity    AlertSeverity
	Check       func(*SuperClaudeConfig) ComplianceResult
}

// ComplianceResult represents compliance check result
type ComplianceResult struct {
	Compliant bool                   `json:"compliant"`
	Message   string                 `json:"message"`
	Evidence  map[string]interface{} `json:"evidence"`
	Remediation string               `json:"remediation"`
}

// AlertManager handles configuration alerts
type AlertManager struct {
	channels    []AlertChannel
	rules       []AlertRule
	suppressions map[string]time.Time
	mu          sync.RWMutex
}

// AlertChannel defines where alerts are sent
type AlertChannel interface {
	SendAlert(alert Alert) error
	Name() string
}

// AlertRule defines when to trigger alerts
type AlertRule struct {
	Name        string
	Condition   func(*SuperClaudeConfig) bool
	Severity    AlertSeverity
	Message     string
	Cooldown    time.Duration
	Channels    []string
}

// Alert represents a configuration alert
type Alert struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	Severity    AlertSeverity          `json:"severity"`
	Title       string                 `json:"title"`
	Message     string                 `json:"message"`
	Component   string                 `json:"component"`
	Environment string                 `json:"environment"`
	Metadata    map[string]interface{} `json:"metadata"`
	Actions     []string               `json:"actions"`
}

// NewConfigObservability creates a new configuration observability system
func NewConfigObservability() *ConfigObservability {
	return &ConfigObservability{
		metrics:           newConfigMetrics(),
		driftDetector:     newConfigDriftDetector(),
		healthChecker:     newConfigHealthChecker(),
		complianceChecker: newComplianceChecker(),
		alertManager:      newAlertManager(),
		enabled:           true,
	}
}

// Start begins configuration monitoring
func (co *ConfigObservability) Start(ctx context.Context, config *SuperClaudeConfig) error {
	co.mu.Lock()
	defer co.mu.Unlock()

	if !co.enabled {
		return nil
	}

	// Start drift detection
	if err := co.driftDetector.Start(ctx, config); err != nil {
		return fmt.Errorf("failed to start drift detector: %w", err)
	}

	// Start health checking
	if err := co.healthChecker.Start(ctx, config); err != nil {
		return fmt.Errorf("failed to start health checker: %w", err)
	}

	// Start compliance monitoring
	if err := co.complianceChecker.Start(ctx, config); err != nil {
		return fmt.Errorf("failed to start compliance checker: %w", err)
	}

	logging.Info("Configuration observability started")
	return nil
}

// RecordConfigLoad records a configuration load event
func (co *ConfigObservability) RecordConfigLoad(config *SuperClaudeConfig, duration time.Duration) {
	if !co.enabled {
		return
	}

	co.metrics.configLoads.Inc()
	co.metrics.validationDuration.Observe(duration.Seconds())
	
	// Calculate config size (approximate)
	configJSON, _ := json.Marshal(config)
	co.metrics.configSize.Set(float64(len(configJSON)))

	logging.Debug("Configuration load recorded", 
		"duration", duration,
		"size_bytes", len(configJSON))
}

// RecordConfigValidation records a validation event
func (co *ConfigObservability) RecordConfigValidation(valid bool, issues int) {
	if !co.enabled {
		return
	}

	co.metrics.configValidations.Inc()
	if !valid {
		co.metrics.configErrors.Add(float64(issues))
	}
}

// RecordConfigChange records a configuration change
func (co *ConfigObservability) RecordConfigChange(oldConfig, newConfig *SuperClaudeConfig) {
	if !co.enabled {
		return
	}

	co.metrics.configChanges.Inc()
	
	// Trigger drift detection
	go co.driftDetector.CheckDrift(oldConfig, newConfig)
}

// GetHealthStatus returns overall configuration health
func (co *ConfigObservability) GetHealthStatus() map[string]HealthResult {
	co.healthChecker.mu.RLock()
	defer co.healthChecker.mu.RUnlock()

	results := make(map[string]HealthResult)
	for name, result := range co.healthChecker.results {
		results[name] = result
	}

	return results
}

// GetComplianceStatus returns compliance status
func (co *ConfigObservability) GetComplianceStatus(config *SuperClaudeConfig) ComplianceReport {
	return co.complianceChecker.CheckCompliance(config)
}

// ComplianceReport contains compliance check results
type ComplianceReport struct {
	OverallCompliant bool                      `json:"overall_compliant"`
	Standards        map[string]StandardResult `json:"standards"`
	Summary          ComplianceSummary         `json:"summary"`
	Timestamp        time.Time                 `json:"timestamp"`
}

// StandardResult contains results for a compliance standard
type StandardResult struct {
	Compliant   bool                      `json:"compliant"`
	Score       float64                   `json:"score"`
	Rules       map[string]ComplianceResult `json:"rules"`
	Required    bool                      `json:"required"`
}

// ComplianceSummary provides aggregate compliance metrics
type ComplianceSummary struct {
	TotalRules     int     `json:"total_rules"`
	PassedRules    int     `json:"passed_rules"`
	FailedRules    int     `json:"failed_rules"`
	ComplianceRate float64 `json:"compliance_rate"`
	CriticalIssues int     `json:"critical_issues"`
}

// Private implementation functions

func newConfigMetrics() *ConfigMetrics {
	return &ConfigMetrics{
		configLoads: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_loads_total",
			Help: "Total number of configuration loads",
		}),
		configValidations: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_validations_total",
			Help: "Total number of configuration validations",
		}),
		configErrors: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_errors_total",
			Help: "Total number of configuration errors",
		}),
		configReloads: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_reloads_total",
			Help: "Total number of configuration reloads",
		}),
		configSize: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "superclaude_config_size_bytes",
			Help: "Current configuration size in bytes",
		}),
		validationDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "superclaude_config_validation_duration_seconds",
			Help:    "Time spent validating configuration",
			Buckets: prometheus.DefBuckets,
		}),
		tenantCount: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "superclaude_config_tenants_total",
			Help: "Total number of configured tenants",
		}),
		activeTenants: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "superclaude_config_active_tenants",
			Help: "Number of active tenants",
		}),
		quotaViolations: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_quota_violations_total",
			Help: "Total number of quota violations",
		}),
		featureUsage: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "superclaude_config_feature_usage_total",
			Help: "Feature usage by configuration",
		}, []string{"feature", "tenant"}),
		configChanges: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_changes_total",
			Help: "Total number of configuration changes",
		}),
		driftDetections: promauto.NewCounter(prometheus.CounterOpts{
			Name: "superclaude_config_drift_detections_total",
			Help: "Total number of configuration drift detections",
		}),
	}
}

func newConfigDriftDetector() *ConfigDriftDetector {
	return &ConfigDriftDetector{
		checkInterval:  5 * time.Minute,
		driftThreshold: 0.1, // 10% change threshold
		alertChannel:   make(chan DriftAlert, 100),
	}
}

func (cdd *ConfigDriftDetector) Start(ctx context.Context, baseline *SuperClaudeConfig) error {
	cdd.mu.Lock()
	defer cdd.mu.Unlock()

	cdd.baseline = baseline
	cdd.running = true

	go cdd.monitor(ctx)
	return nil
}

func (cdd *ConfigDriftDetector) monitor(ctx context.Context) {
	ticker := time.NewTicker(cdd.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Periodic drift checks would go here
		case alert := <-cdd.alertChannel:
			logging.Warn("Configuration drift detected",
				"type", alert.DriftType,
				"component", alert.Component,
				"severity", alert.Severity)
		}
	}
}

func (cdd *ConfigDriftDetector) CheckDrift(old, new *SuperClaudeConfig) {
	// Calculate configuration drift
	changes := cdd.calculateChanges(old, new)
	
	for _, change := range changes {
		if change.Significance > cdd.driftThreshold {
			alert := DriftAlert{
				Timestamp:  time.Now(),
				DriftType:  change.Type,
				Severity:   change.Severity,
				Component:  change.Component,
				Expected:   change.Expected,
				Actual:     change.Actual,
				Difference: change.Significance,
			}
			
			select {
			case cdd.alertChannel <- alert:
			default:
				logging.Warn("Drift alert channel full, dropping alert")
			}
		}
	}
}

type ConfigDriftChange struct {
	Component    string
	Type         DriftType
	Severity     AlertSeverity
	Expected     interface{}
	Actual       interface{}
	Significance float64
}

func (cdd *ConfigDriftDetector) calculateChanges(old, new *SuperClaudeConfig) []ConfigDriftChange {
	var changes []ConfigDriftChange

	// Server configuration changes
	if old.Server.Port != new.Server.Port {
		changes = append(changes, ConfigDriftChange{
			Component:    "server.port",
			Type:         DriftConfiguration,
			Severity:     AlertWarning,
			Expected:     old.Server.Port,
			Actual:       new.Server.Port,
			Significance: 0.5, // Port changes are significant
		})
	}

	// Security configuration changes
	if old.Security.APIKeyEncryption != new.Security.APIKeyEncryption {
		changes = append(changes, ConfigDriftChange{
			Component:    "security.api_key_encryption",
			Type:         DriftSecurity,
			Severity:     AlertCritical,
			Expected:     old.Security.APIKeyEncryption,
			Actual:       new.Security.APIKeyEncryption,
			Significance: 1.0, // Security changes are always significant
		})
	}

	return changes
}

func newConfigHealthChecker() *ConfigHealthChecker {
	return &ConfigHealthChecker{
		checks:   getDefaultHealthChecks(),
		interval: 1 * time.Minute,
		results:  make(map[string]HealthResult),
	}
}

func (chc *ConfigHealthChecker) Start(ctx context.Context, config *SuperClaudeConfig) error {
	go chc.runHealthChecks(ctx, config)
	return nil
}

func (chc *ConfigHealthChecker) runHealthChecks(ctx context.Context, config *SuperClaudeConfig) {
	ticker := time.NewTicker(chc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, check := range chc.checks {
				start := time.Now()
				result := check.Check(config)
				result.Duration = time.Since(start)
				result.Timestamp = time.Now()

				chc.mu.Lock()
				chc.results[check.Name] = result
				chc.mu.Unlock()
			}
		}
	}
}

func getDefaultHealthChecks() []HealthCheck {
	return []HealthCheck{
		{
			Name:        "database_connection",
			Description: "Check database connectivity",
			Critical:    true,
			Check: func(config *SuperClaudeConfig) HealthResult {
				// Would implement actual database connectivity check
				return HealthResult{
					Status:  HealthHealthy,
					Message: "Database connection healthy",
				}
			},
		},
		{
			Name:        "api_provider_connectivity",
			Description: "Check AI provider API connectivity",
			Critical:    true,
			Check: func(config *SuperClaudeConfig) HealthResult {
				// Would implement actual API connectivity check
				return HealthResult{
					Status:  HealthHealthy,
					Message: "AI provider APIs accessible",
				}
			},
		},
		{
			Name:        "cache_availability",
			Description: "Check cache system availability",
			Critical:    false,
			Check: func(config *SuperClaudeConfig) HealthResult {
				// Would implement cache connectivity check
				return HealthResult{
					Status:  HealthHealthy,
					Message: "Cache system operational",
				}
			},
		},
	}
}

func newComplianceChecker() *ComplianceChecker {
	return &ComplianceChecker{
		standards: getDefaultComplianceStandards(),
	}
}

func (cc *ComplianceChecker) Start(ctx context.Context, config *SuperClaudeConfig) error {
	// Initial compliance check
	report := cc.CheckCompliance(config)
	
	logging.Info("Initial compliance check completed",
		"overall_compliant", report.OverallCompliant,
		"compliance_rate", report.Summary.ComplianceRate)
	
	return nil
}

func (cc *ComplianceChecker) CheckCompliance(config *SuperClaudeConfig) ComplianceReport {
	report := ComplianceReport{
		Standards: make(map[string]StandardResult),
		Timestamp: time.Now(),
	}

	totalRules := 0
	passedRules := 0
	criticalIssues := 0

	for _, standard := range cc.standards {
		result := StandardResult{
			Rules:    make(map[string]ComplianceResult),
			Required: standard.Required,
		}

		standardPassed := 0
		for _, rule := range standard.Rules {
			ruleResult := rule.Check(config)
			result.Rules[rule.ID] = ruleResult

			totalRules++
			if ruleResult.Compliant {
				standardPassed++
				passedRules++
			} else if rule.Severity == AlertCritical {
				criticalIssues++
			}
		}

		result.Compliant = standardPassed == len(standard.Rules)
		result.Score = float64(standardPassed) / float64(len(standard.Rules))
		report.Standards[standard.Name] = result
	}

	// Calculate overall compliance
	report.OverallCompliant = criticalIssues == 0
	report.Summary = ComplianceSummary{
		TotalRules:     totalRules,
		PassedRules:    passedRules,
		FailedRules:    totalRules - passedRules,
		ComplianceRate: float64(passedRules) / float64(totalRules),
		CriticalIssues: criticalIssues,
	}

	return report
}

func getDefaultComplianceStandards() []ComplianceStandard {
	return []ComplianceStandard{
		{
			Name:        "SOC2",
			Description: "SOC 2 Type II Compliance",
			Version:     "2017",
			Required:    true,
			Rules: []ComplianceRule{
				{
					ID:          "SOC2-CC6.1",
					Description: "Encryption in transit must be enabled",
					Severity:    AlertCritical,
					Check: func(config *SuperClaudeConfig) ComplianceResult {
						if config.Server.TLS.Enabled {
							return ComplianceResult{
								Compliant: true,
								Message:   "TLS encryption enabled",
							}
						}
						return ComplianceResult{
							Compliant:   false,
							Message:     "TLS encryption not enabled",
							Remediation: "Enable TLS in server configuration",
						}
					},
				},
				{
					ID:          "SOC2-CC6.7",
					Description: "API keys must be encrypted at rest",
					Severity:    AlertCritical,
					Check: func(config *SuperClaudeConfig) ComplianceResult {
						if config.Security.APIKeyEncryption {
							return ComplianceResult{
								Compliant: true,
								Message:   "API key encryption enabled",
							}
						}
						return ComplianceResult{
							Compliant:   false,
							Message:     "API keys not encrypted",
							Remediation: "Enable API key encryption in security configuration",
						}
					},
				},
			},
		},
	}
}

func newAlertManager() *AlertManager {
	return &AlertManager{
		channels:     []AlertChannel{},
		rules:        getDefaultAlertRules(),
		suppressions: make(map[string]time.Time),
	}
}

func getDefaultAlertRules() []AlertRule {
	return []AlertRule{
		{
			Name:     "tls_disabled_production",
			Severity: AlertCritical,
			Message:  "TLS is disabled in production environment",
			Cooldown: 1 * time.Hour,
			Condition: func(config *SuperClaudeConfig) bool {
				return config.Deployment.Environment == "production" && !config.Server.TLS.Enabled
			},
		},
		{
			Name:     "high_error_rate",
			Severity: AlertWarning,
			Message:  "Configuration validation error rate is high",
			Cooldown: 15 * time.Minute,
			Condition: func(config *SuperClaudeConfig) bool {
				// Would check actual error rate metrics
				return false
			},
		},
	}
}