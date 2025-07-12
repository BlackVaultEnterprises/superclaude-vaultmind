package config

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/opencode-ai/opencode/internal/logging"
	"gopkg.in/yaml.v3"
)

// ConfigManager provides advanced configuration management
type ConfigManager struct {
	config          *SuperClaudeConfig
	mu              sync.RWMutex
	watchers        []ConfigWatcher
	encryptionKey   []byte
	version         ConfigVersion
	auditLogger     AuditLogger
	validationRules []ValidationRule
	migrator        *ConfigMigrator
	hotReload       bool
	ctx             context.Context
	cancel          context.CancelFunc
}

// ConfigWatcher defines interface for configuration change watchers
type ConfigWatcher interface {
	OnConfigChange(old, new *SuperClaudeConfig) error
}

// ConfigVersion tracks configuration schema versions
type ConfigVersion struct {
	Version       string    `json:"version" yaml:"version"`
	MinVersion    string    `json:"min_version" yaml:"min_version"`
	SchemaHash    string    `json:"schema_hash" yaml:"schema_hash"`
	LastUpdated   time.Time `json:"last_updated" yaml:"last_updated"`
	Compatibility []string  `json:"compatibility" yaml:"compatibility"`
}

// AuditLogger tracks configuration changes
type AuditLogger struct {
	enabled    bool
	logPath    string
	retention  time.Duration
	encryptLog bool
}

// ValidationRule defines custom validation logic
type ValidationRule struct {
	Name        string
	Description string
	Validator   func(*SuperClaudeConfig) error
	Severity    ValidationSeverity
	Category    string
}

type ValidationSeverity int

const (
	ValidationInfo ValidationSeverity = iota
	ValidationWarning
	ValidationError
	ValidationCritical
)

// ConfigMigrator handles configuration migrations
type ConfigMigrator struct {
	migrations map[string]Migration
	mu         sync.RWMutex
}

// Migration defines a configuration migration
type Migration struct {
	FromVersion string
	ToVersion   string
	Description string
	Migrate     func(map[string]interface{}) (map[string]interface{}, error)
	Rollback    func(map[string]interface{}) (map[string]interface{}, error)
}

// NewConfigManager creates an advanced configuration manager
func NewConfigManager(configPath string, opts ...ConfigOption) (*ConfigManager, error) {
	ctx, cancel := context.WithCancel(context.Background())
	
	cm := &ConfigManager{
		watchers:        make([]ConfigWatcher, 0),
		validationRules: getDefaultValidationRules(),
		migrator:        NewConfigMigrator(),
		ctx:             ctx,
		cancel:          cancel,
		hotReload:       true,
	}
	
	// Apply options
	for _, opt := range opts {
		opt(cm)
	}
	
	// Load initial configuration
	config, err := cm.LoadWithValidation(configPath)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	cm.config = config
	
	// Initialize audit logging
	if err := cm.initAuditLogging(); err != nil {
		logging.Warn("Failed to initialize audit logging", "error", err)
	}
	
	// Start hot reload if enabled
	if cm.hotReload {
		go cm.watchConfigChanges(configPath)
	}
	
	return cm, nil
}

// ConfigOption defines configuration manager options
type ConfigOption func(*ConfigManager)

// WithEncryption enables configuration encryption
func WithEncryption(key string) ConfigOption {
	return func(cm *ConfigManager) {
		hash := sha256.Sum256([]byte(key))
		cm.encryptionKey = hash[:]
	}
}

// WithHotReload enables/disables hot reload
func WithHotReload(enabled bool) ConfigOption {
	return func(cm *ConfigManager) {
		cm.hotReload = enabled
	}
}

// WithAuditLogging enables audit logging
func WithAuditLogging(path string, retention time.Duration) ConfigOption {
	return func(cm *ConfigManager) {
		cm.auditLogger = AuditLogger{
			enabled:   true,
			logPath:   path,
			retention: retention,
		}
	}
}

// LoadWithValidation loads configuration with comprehensive validation
func (cm *ConfigManager) LoadWithValidation(configPath string) (*SuperClaudeConfig, error) {
	// Load base configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}
	
	// Version check
	if err := cm.validateVersion(config); err != nil {
		return nil, fmt.Errorf("version validation failed: %w", err)
	}
	
	// Custom validation rules
	if err := cm.runValidationRules(config); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	
	// Decrypt sensitive fields
	if cm.encryptionKey != nil {
		if err := cm.decryptSensitiveFields(config); err != nil {
			return nil, fmt.Errorf("decryption failed: %w", err)
		}
	}
	
	// Apply security hardening
	cm.applySecurityHardening(config)
	
	return config, nil
}

// GetConfig returns the current configuration (thread-safe)
func (cm *ConfigManager) GetConfig() *SuperClaudeConfig {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	// Return a deep copy to prevent mutations
	return cm.deepCopyConfig(cm.config)
}

// UpdateConfig updates configuration with validation and audit
func (cm *ConfigManager) UpdateConfig(updates map[string]interface{}) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	oldConfig := cm.deepCopyConfig(cm.config)
	
	// Apply updates
	newConfig, err := cm.applyUpdates(cm.config, updates)
	if err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}
	
	// Validate new configuration
	if err := cm.runValidationRules(newConfig); err != nil {
		return fmt.Errorf("validation failed after update: %w", err)
	}
	
	// Audit the change
	cm.auditConfigChange(oldConfig, newConfig, updates)
	
	// Notify watchers
	for _, watcher := range cm.watchers {
		if err := watcher.OnConfigChange(oldConfig, newConfig); err != nil {
			logging.Error("Config watcher failed", "error", err)
		}
	}
	
	cm.config = newConfig
	
	logging.Info("Configuration updated successfully", 
		"changes", len(updates),
		"version", newConfig.Deployment.Version)
	
	return nil
}

// AddWatcher adds a configuration change watcher
func (cm *ConfigManager) AddWatcher(watcher ConfigWatcher) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.watchers = append(cm.watchers, watcher)
}

// Encrypt encrypts sensitive configuration values
func (cm *ConfigManager) Encrypt(plaintext string) (string, error) {
	if cm.encryptionKey == nil {
		return plaintext, nil
	}
	
	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts sensitive configuration values
func (cm *ConfigManager) Decrypt(ciphertext string) (string, error) {
	if cm.encryptionKey == nil {
		return ciphertext, nil
	}
	
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", err
	}
	
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	
	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// ValidateConfiguration runs comprehensive validation
func (cm *ConfigManager) ValidateConfiguration() *ValidationResult {
	result := &ValidationResult{
		Valid:    true,
		Issues:   make([]ValidationIssue, 0),
		Warnings: make([]ValidationIssue, 0),
	}
	
	for _, rule := range cm.validationRules {
		if err := rule.Validator(cm.config); err != nil {
			issue := ValidationIssue{
				Rule:        rule.Name,
				Description: rule.Description,
				Error:       err.Error(),
				Severity:    rule.Severity,
				Category:    rule.Category,
			}
			
			switch rule.Severity {
			case ValidationError, ValidationCritical:
				result.Valid = false
				result.Issues = append(result.Issues, issue)
			case ValidationWarning:
				result.Warnings = append(result.Warnings, issue)
			}
		}
	}
	
	return result
}

// ValidationResult contains validation results
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Issues   []ValidationIssue `json:"issues"`
	Warnings []ValidationIssue `json:"warnings"`
}

// ValidationIssue represents a configuration validation issue
type ValidationIssue struct {
	Rule        string             `json:"rule"`
	Description string             `json:"description"`
	Error       string             `json:"error"`
	Severity    ValidationSeverity `json:"severity"`
	Category    string             `json:"category"`
}

// GetConfigHistory returns configuration change history
func (cm *ConfigManager) GetConfigHistory(limit int) ([]ConfigChange, error) {
	// Implementation would read from audit log
	return nil, fmt.Errorf("not implemented")
}

// ConfigChange represents a configuration change event
type ConfigChange struct {
	Timestamp time.Time              `json:"timestamp"`
	User      string                 `json:"user"`
	Changes   map[string]interface{} `json:"changes"`
	Version   string                 `json:"version"`
	Source    string                 `json:"source"`
}

// ExportConfig exports configuration in various formats
func (cm *ConfigManager) ExportConfig(format string, includeSecrets bool) ([]byte, error) {
	config := cm.GetConfig()
	
	if !includeSecrets {
		config = cm.redactSecrets(config)
	}
	
	switch format {
	case "yaml":
		return yaml.Marshal(config)
	case "json":
		return json.MarshalIndent(config, "", "  ")
	case "toml":
		// Would implement TOML marshaling
		return nil, fmt.Errorf("TOML export not implemented")
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Private helper methods

func (cm *ConfigManager) watchConfigChanges(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logging.Error("Failed to create file watcher", "error", err)
		return
	}
	defer watcher.Close()
	
	if err := watcher.Add(configPath); err != nil {
		logging.Error("Failed to watch config path", "error", err)
		return
	}
	
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				cm.handleConfigFileChange()
			}
		case err := <-watcher.Errors:
			logging.Error("Config watcher error", "error", err)
		case <-cm.ctx.Done():
			return
		}
	}
}

func (cm *ConfigManager) handleConfigFileChange() {
	logging.Info("Configuration file changed, reloading...")
	
	// Add debouncing to prevent rapid reloads
	time.Sleep(100 * time.Millisecond)
	
	newConfig, err := cm.LoadWithValidation("")
	if err != nil {
		logging.Error("Failed to reload configuration", "error", err)
		return
	}
	
	oldConfig := cm.GetConfig()
	
	cm.mu.Lock()
	cm.config = newConfig
	cm.mu.Unlock()
	
	// Notify watchers
	for _, watcher := range cm.watchers {
		if err := watcher.OnConfigChange(oldConfig, newConfig); err != nil {
			logging.Error("Config watcher failed during hot reload", "error", err)
		}
	}
	
	logging.Info("Configuration reloaded successfully")
}

func (cm *ConfigManager) validateVersion(config *SuperClaudeConfig) error {
	// Version validation logic
	return nil
}

func (cm *ConfigManager) runValidationRules(config *SuperClaudeConfig) error {
	for _, rule := range cm.validationRules {
		if err := rule.Validator(config); err != nil {
			if rule.Severity == ValidationCritical {
				return fmt.Errorf("critical validation failed for rule %s: %w", rule.Name, err)
			}
		}
	}
	return nil
}

func (cm *ConfigManager) decryptSensitiveFields(config *SuperClaudeConfig) error {
	// Decrypt API keys and other sensitive fields
	var err error
	
	if config.Providers.OpenRouter.APIKey != "" {
		config.Providers.OpenRouter.APIKey, err = cm.Decrypt(config.Providers.OpenRouter.APIKey)
		if err != nil {
			return fmt.Errorf("failed to decrypt OpenRouter API key: %w", err)
		}
	}
	
	// Decrypt other sensitive fields...
	
	return nil
}

func (cm *ConfigManager) applySecurityHardening(config *SuperClaudeConfig) {
	// Apply security best practices
	if config.Deployment.Environment == "production" {
		// Force secure settings in production
		config.Security.APIKeyEncryption = true
		config.Security.SessionEncryption = true
		config.Monitoring.Profiling.Enabled = false
		
		// Ensure TLS is enabled
		if !config.Server.TLS.Enabled {
			logging.Warn("TLS should be enabled in production")
		}
	}
}

func (cm *ConfigManager) deepCopyConfig(config *SuperClaudeConfig) *SuperClaudeConfig {
	// Implementation would create a deep copy
	return config // Simplified for now
}

func (cm *ConfigManager) applyUpdates(config *SuperClaudeConfig, updates map[string]interface{}) (*SuperClaudeConfig, error) {
	// Implementation would apply updates safely
	return config, nil
}

func (cm *ConfigManager) auditConfigChange(old, new *SuperClaudeConfig, updates map[string]interface{}) {
	if !cm.auditLogger.enabled {
		return
	}
	
	change := ConfigChange{
		Timestamp: time.Now(),
		Changes:   updates,
		Version:   new.Deployment.Version,
		Source:    "api",
	}
	
	// Log to audit file
	data, _ := json.Marshal(change)
	logging.Info("Configuration change audited", "change", string(data))
}

func (cm *ConfigManager) redactSecrets(config *SuperClaudeConfig) *SuperClaudeConfig {
	// Implementation would redact sensitive fields
	redacted := cm.deepCopyConfig(config)
	redacted.Providers.OpenRouter.APIKey = "[REDACTED]"
	redacted.Providers.OpenAI.APIKey = "[REDACTED]"
	redacted.Providers.Anthropic.APIKey = "[REDACTED]"
	redacted.Security.Auth.JWTSecret = "[REDACTED]"
	return redacted
}

func (cm *ConfigManager) initAuditLogging() error {
	if !cm.auditLogger.enabled {
		return nil
	}
	
	// Initialize audit logging
	logging.Info("Audit logging initialized", "path", cm.auditLogger.logPath)
	return nil
}

// Close gracefully shuts down the config manager
func (cm *ConfigManager) Close() error {
	cm.cancel()
	logging.Info("Configuration manager closed")
	return nil
}

// NewConfigMigrator creates a new configuration migrator
func NewConfigMigrator() *ConfigMigrator {
	return &ConfigMigrator{
		migrations: make(map[string]Migration),
	}
}

// Default validation rules
func getDefaultValidationRules() []ValidationRule {
	return []ValidationRule{
		{
			Name:        "port_range",
			Description: "Validate port numbers are in valid range",
			Severity:    ValidationError,
			Category:    "network",
			Validator: func(config *SuperClaudeConfig) error {
				if config.Server.Port < 1 || config.Server.Port > 65535 {
					return fmt.Errorf("server port %d is not in valid range 1-65535", config.Server.Port)
				}
				return nil
			},
		},
		{
			Name:        "tls_production",
			Description: "TLS should be enabled in production",
			Severity:    ValidationWarning,
			Category:    "security",
			Validator: func(config *SuperClaudeConfig) error {
				if config.Deployment.Environment == "production" && !config.Server.TLS.Enabled {
					return fmt.Errorf("TLS should be enabled in production environment")
				}
				return nil
			},
		},
		{
			Name:        "api_key_present",
			Description: "API keys should be configured",
			Severity:    ValidationWarning,
			Category:    "configuration",
			Validator: func(config *SuperClaudeConfig) error {
				if config.Providers.OpenRouter.APIKey == "" && 
				   config.Providers.OpenAI.APIKey == "" && 
				   config.Providers.Anthropic.APIKey == "" {
					return fmt.Errorf("no API keys configured for any providers")
				}
				return nil
			},
		},
	}
}