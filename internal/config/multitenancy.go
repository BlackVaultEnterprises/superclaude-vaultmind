package config

import (
	"fmt"
	"sync"
	"time"
)

// MultiTenantConfigManager manages configurations for multiple tenants
type MultiTenantConfigManager struct {
	tenants       map[string]*TenantConfig
	globalConfig  *SuperClaudeConfig
	mu            sync.RWMutex
	defaultTenant string
	isolation     IsolationLevel
}

// TenantConfig represents tenant-specific configuration
type TenantConfig struct {
	ID            string                 `json:"id" yaml:"id"`
	Name          string                 `json:"name" yaml:"name"`
	Config        *SuperClaudeConfig     `json:"config" yaml:"config"`
	Overrides     map[string]interface{} `json:"overrides" yaml:"overrides"`
	Quotas        *TenantQuotas          `json:"quotas" yaml:"quotas"`
	Features      *TenantFeatures        `json:"features" yaml:"features"`
	Metadata      map[string]string      `json:"metadata" yaml:"metadata"`
	CreatedAt     time.Time              `json:"created_at" yaml:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at" yaml:"updated_at"`
	Status        TenantStatus           `json:"status" yaml:"status"`
}

// TenantQuotas defines resource quotas per tenant
type TenantQuotas struct {
	MaxSessions           int           `json:"max_sessions" yaml:"max_sessions"`
	MaxRequestsPerMinute  int           `json:"max_requests_per_minute" yaml:"max_requests_per_minute"`
	MaxTokensPerMonth     int64         `json:"max_tokens_per_month" yaml:"max_tokens_per_month"`
	MaxStorageMB          int           `json:"max_storage_mb" yaml:"max_storage_mb"`
	MaxConcurrentRequests int           `json:"max_concurrent_requests" yaml:"max_concurrent_requests"`
	SessionTimeout        time.Duration `json:"session_timeout" yaml:"session_timeout"`
	DataRetention         time.Duration `json:"data_retention" yaml:"data_retention"`
}

// TenantFeatures defines which features are enabled per tenant
type TenantFeatures struct {
	MCPServer         bool `json:"mcp_server" yaml:"mcp_server"`
	AdvancedPersonas  bool `json:"advanced_personas" yaml:"advanced_personas"`
	CustomCommands    bool `json:"custom_commands" yaml:"custom_commands"`
	APIAccess         bool `json:"api_access" yaml:"api_access"`
	AuditLogging      bool `json:"audit_logging" yaml:"audit_logging"`
	PrioritySupport   bool `json:"priority_support" yaml:"priority_support"`
	CustomIntegration bool `json:"custom_integration" yaml:"custom_integration"`
	AdvancedAnalytics bool `json:"advanced_analytics" yaml:"advanced_analytics"`
}

type TenantStatus int

const (
	TenantActive TenantStatus = iota
	TenantSuspended
	TenantDeactivated
	TenantMaintenance
)

type IsolationLevel int

const (
	IsolationShared IsolationLevel = iota
	IsolationDedicated
	IsolationPrivate
)

// NewMultiTenantConfigManager creates a new multi-tenant configuration manager
func NewMultiTenantConfigManager(globalConfig *SuperClaudeConfig, isolation IsolationLevel) *MultiTenantConfigManager {
	return &MultiTenantConfigManager{
		tenants:       make(map[string]*TenantConfig),
		globalConfig:  globalConfig,
		isolation:     isolation,
		defaultTenant: "default",
	}
}

// CreateTenant creates a new tenant configuration
func (mtcm *MultiTenantConfigManager) CreateTenant(tenantID, name string, quotas *TenantQuotas, features *TenantFeatures) (*TenantConfig, error) {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	if _, exists := mtcm.tenants[tenantID]; exists {
		return nil, fmt.Errorf("tenant %s already exists", tenantID)
	}
	
	// Create tenant-specific config based on global config
	tenantConfig := mtcm.createTenantConfig(tenantID, name, quotas, features)
	
	mtcm.tenants[tenantID] = tenantConfig
	
	return tenantConfig, nil
}

// GetTenantConfig returns configuration for a specific tenant
func (mtcm *MultiTenantConfigManager) GetTenantConfig(tenantID string) (*SuperClaudeConfig, error) {
	mtcm.mu.RLock()
	defer mtcm.mu.RUnlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		if tenantID == mtcm.defaultTenant {
			return mtcm.globalConfig, nil
		}
		return nil, fmt.Errorf("tenant %s not found", tenantID)
	}
	
	if tenant.Status != TenantActive {
		return nil, fmt.Errorf("tenant %s is not active (status: %v)", tenantID, tenant.Status)
	}
	
	return tenant.Config, nil
}

// UpdateTenantConfig updates configuration for a specific tenant
func (mtcm *MultiTenantConfigManager) UpdateTenantConfig(tenantID string, overrides map[string]interface{}) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}
	
	// Apply overrides to tenant config
	if err := mtcm.applyTenantOverrides(tenant, overrides); err != nil {
		return fmt.Errorf("failed to apply overrides: %w", err)
	}
	
	tenant.UpdatedAt = time.Now()
	
	return nil
}

// DeleteTenant removes a tenant configuration
func (mtcm *MultiTenantConfigManager) DeleteTenant(tenantID string) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	if tenantID == mtcm.defaultTenant {
		return fmt.Errorf("cannot delete default tenant")
	}
	
	delete(mtcm.tenants, tenantID)
	return nil
}

// ListTenants returns all tenant configurations
func (mtcm *MultiTenantConfigManager) ListTenants() []*TenantConfig {
	mtcm.mu.RLock()
	defer mtcm.mu.RUnlock()
	
	tenants := make([]*TenantConfig, 0, len(mtcm.tenants))
	for _, tenant := range mtcm.tenants {
		tenants = append(tenants, tenant)
	}
	
	return tenants
}

// SetTenantStatus updates tenant status
func (mtcm *MultiTenantConfigManager) SetTenantStatus(tenantID string, status TenantStatus) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}
	
	tenant.Status = status
	tenant.UpdatedAt = time.Now()
	
	return nil
}

// ValidateTenantQuotas validates that tenant usage is within quotas
func (mtcm *MultiTenantConfigManager) ValidateTenantQuotas(tenantID string, usage *TenantUsage) error {
	mtcm.mu.RLock()
	defer mtcm.mu.RUnlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}
	
	quotas := tenant.Quotas
	if quotas == nil {
		return nil // No quotas defined
	}
	
	// Check various quota limits
	if usage.ActiveSessions > quotas.MaxSessions {
		return fmt.Errorf("tenant %s exceeded max sessions: %d > %d", 
			tenantID, usage.ActiveSessions, quotas.MaxSessions)
	}
	
	if usage.RequestsPerMinute > quotas.MaxRequestsPerMinute {
		return fmt.Errorf("tenant %s exceeded requests per minute: %d > %d", 
			tenantID, usage.RequestsPerMinute, quotas.MaxRequestsPerMinute)
	}
	
	if usage.TokensThisMonth > quotas.MaxTokensPerMonth {
		return fmt.Errorf("tenant %s exceeded monthly token limit: %d > %d", 
			tenantID, usage.TokensThisMonth, quotas.MaxTokensPerMonth)
	}
	
	if usage.StorageUsedMB > quotas.MaxStorageMB {
		return fmt.Errorf("tenant %s exceeded storage limit: %d MB > %d MB", 
			tenantID, usage.StorageUsedMB, quotas.MaxStorageMB)
	}
	
	return nil
}

// TenantUsage represents current usage metrics for a tenant
type TenantUsage struct {
	ActiveSessions      int   `json:"active_sessions"`
	RequestsPerMinute   int   `json:"requests_per_minute"`
	TokensThisMonth     int64 `json:"tokens_this_month"`
	StorageUsedMB       int   `json:"storage_used_mb"`
	ConcurrentRequests  int   `json:"concurrent_requests"`
	LastActivity        time.Time `json:"last_activity"`
}

// GetTenantUsage returns current usage for a tenant
func (mtcm *MultiTenantConfigManager) GetTenantUsage(tenantID string) (*TenantUsage, error) {
	// This would integrate with metrics/monitoring systems
	// For now, return empty usage
	return &TenantUsage{}, nil
}

// EnableFeatureForTenant enables a specific feature for a tenant
func (mtcm *MultiTenantConfigManager) EnableFeatureForTenant(tenantID string, feature string) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}
	
	if tenant.Features == nil {
		tenant.Features = &TenantFeatures{}
	}
	
	switch feature {
	case "mcp_server":
		tenant.Features.MCPServer = true
	case "advanced_personas":
		tenant.Features.AdvancedPersonas = true
	case "custom_commands":
		tenant.Features.CustomCommands = true
	case "api_access":
		tenant.Features.APIAccess = true
	case "audit_logging":
		tenant.Features.AuditLogging = true
	case "priority_support":
		tenant.Features.PrioritySupport = true
	case "custom_integration":
		tenant.Features.CustomIntegration = true
	case "advanced_analytics":
		tenant.Features.AdvancedAnalytics = true
	default:
		return fmt.Errorf("unknown feature: %s", feature)
	}
	
	tenant.UpdatedAt = time.Now()
	return nil
}

// GetTenantsByFeature returns all tenants with a specific feature enabled
func (mtcm *MultiTenantConfigManager) GetTenantsByFeature(feature string) []*TenantConfig {
	mtcm.mu.RLock()
	defer mtcm.mu.RUnlock()
	
	var tenants []*TenantConfig
	
	for _, tenant := range mtcm.tenants {
		if tenant.Features == nil {
			continue
		}
		
		var hasFeature bool
		switch feature {
		case "mcp_server":
			hasFeature = tenant.Features.MCPServer
		case "advanced_personas":
			hasFeature = tenant.Features.AdvancedPersonas
		case "custom_commands":
			hasFeature = tenant.Features.CustomCommands
		case "api_access":
			hasFeature = tenant.Features.APIAccess
		case "audit_logging":
			hasFeature = tenant.Features.AuditLogging
		case "priority_support":
			hasFeature = tenant.Features.PrioritySupport
		case "custom_integration":
			hasFeature = tenant.Features.CustomIntegration
		case "advanced_analytics":
			hasFeature = tenant.Features.AdvancedAnalytics
		}
		
		if hasFeature {
			tenants = append(tenants, tenant)
		}
	}
	
	return tenants
}

// ArchiveTenant archives a tenant's data and configuration
func (mtcm *MultiTenantConfigManager) ArchiveTenant(tenantID string) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	tenant, exists := mtcm.tenants[tenantID]
	if !exists {
		return fmt.Errorf("tenant %s not found", tenantID)
	}
	
	// Set tenant to deactivated
	tenant.Status = TenantDeactivated
	tenant.UpdatedAt = time.Now()
	
	// Here you would implement actual archival logic:
	// - Export tenant data
	// - Remove from active systems
	// - Store in archive storage
	
	return nil
}

// BulkUpdateTenants applies updates to multiple tenants
func (mtcm *MultiTenantConfigManager) BulkUpdateTenants(tenantIDs []string, updates map[string]interface{}) error {
	mtcm.mu.Lock()
	defer mtcm.mu.Unlock()
	
	var errors []error
	
	for _, tenantID := range tenantIDs {
		tenant, exists := mtcm.tenants[tenantID]
		if !exists {
			errors = append(errors, fmt.Errorf("tenant %s not found", tenantID))
			continue
		}
		
		if err := mtcm.applyTenantOverrides(tenant, updates); err != nil {
			errors = append(errors, fmt.Errorf("failed to update tenant %s: %w", tenantID, err))
			continue
		}
		
		tenant.UpdatedAt = time.Now()
	}
	
	if len(errors) > 0 {
		return fmt.Errorf("bulk update failed for some tenants: %v", errors)
	}
	
	return nil
}

// Private helper methods

func (mtcm *MultiTenantConfigManager) createTenantConfig(tenantID, name string, quotas *TenantQuotas, features *TenantFeatures) *TenantConfig {
	// Deep copy global config for tenant
	tenantConfig := mtcm.deepCopyConfig(mtcm.globalConfig)
	
	// Apply tenant-specific modifications
	mtcm.applyTenantIsolation(tenantConfig, tenantID)
	
	// Set default quotas if none provided
	if quotas == nil {
		quotas = mtcm.getDefaultQuotas()
	}
	
	// Set default features if none provided
	if features == nil {
		features = mtcm.getDefaultFeatures()
	}
	
	return &TenantConfig{
		ID:        tenantID,
		Name:      name,
		Config:    tenantConfig,
		Overrides: make(map[string]interface{}),
		Quotas:    quotas,
		Features:  features,
		Metadata:  make(map[string]string),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status:    TenantActive,
	}
}

func (mtcm *MultiTenantConfigManager) applyTenantIsolation(config *SuperClaudeConfig, tenantID string) {
	switch mtcm.isolation {
	case IsolationDedicated:
		// Each tenant gets dedicated resources
		config.Database.SQLite.Path = fmt.Sprintf("~/.superclaude/tenants/%s/data.db", tenantID)
		config.Logging.File.Path = fmt.Sprintf("~/.superclaude/tenants/%s/logs/", tenantID)
		
	case IsolationPrivate:
		// Complete isolation with separate infrastructure
		config.Server.Port = config.Server.Port + hashTenantID(tenantID)%1000
		config.MCP.Port = config.MCP.Port + hashTenantID(tenantID)%1000
		
	case IsolationShared:
		// Shared infrastructure with logical separation
		config.Logging.StructuredFields["tenant_id"] = tenantID
	}
}

func (mtcm *MultiTenantConfigManager) applyTenantOverrides(tenant *TenantConfig, overrides map[string]interface{}) error {
	// Apply overrides to tenant configuration
	// This would use reflection or a configuration library to apply nested updates
	
	for key, value := range overrides {
		tenant.Overrides[key] = value
	}
	
	return nil
}

func (mtcm *MultiTenantConfigManager) deepCopyConfig(config *SuperClaudeConfig) *SuperClaudeConfig {
	// Implementation would create a deep copy
	// For now, return the original (this should be implemented properly)
	return config
}

func (mtcm *MultiTenantConfigManager) getDefaultQuotas() *TenantQuotas {
	return &TenantQuotas{
		MaxSessions:           10,
		MaxRequestsPerMinute:  100,
		MaxTokensPerMonth:     1000000,
		MaxStorageMB:          1024,
		MaxConcurrentRequests: 5,
		SessionTimeout:        24 * time.Hour,
		DataRetention:         30 * 24 * time.Hour,
	}
}

func (mtcm *MultiTenantConfigManager) getDefaultFeatures() *TenantFeatures {
	return &TenantFeatures{
		MCPServer:         true,
		AdvancedPersonas:  false,
		CustomCommands:    false,
		APIAccess:         true,
		AuditLogging:      false,
		PrioritySupport:   false,
		CustomIntegration: false,
		AdvancedAnalytics: false,
	}
}

func hashTenantID(tenantID string) int {
	hash := 0
	for _, char := range tenantID {
		hash = (hash*31 + int(char)) % 1000
	}
	return hash
}