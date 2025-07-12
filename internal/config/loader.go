package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// SuperClaudeConfig represents the complete configuration
type SuperClaudeConfig struct {
	Server      ServerConfig      `mapstructure:"server"`
	MCP         MCPConfig         `mapstructure:"mcp"`
	Providers   ProvidersConfig   `mapstructure:"providers"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Cache       CacheConfig       `mapstructure:"cache"`
	Performance PerformanceConfig `mapstructure:"performance"`
	RateLimit   RateLimitConfig   `mapstructure:"rate_limiting"`
	Security    SecurityConfig    `mapstructure:"security"`
	Logging     LoggingConfig     `mapstructure:"logging"`
	Monitoring  MonitoringConfig  `mapstructure:"monitoring"`
	SuperClaude SuperClaudeSpecificConfig `mapstructure:"superclaude"`
	IDE         IDEConfig         `mapstructure:"ide"`
	Features    FeaturesConfig    `mapstructure:"features"`
	Development DevelopmentConfig `mapstructure:"development"`
	Deployment  DeploymentConfig  `mapstructure:"deployment"`
}

type ServerConfig struct {
	Host           string        `mapstructure:"host"`
	Port           int           `mapstructure:"port"`
	Timeout        time.Duration `mapstructure:"timeout"`
	MaxConnections int           `mapstructure:"max_connections"`
	TLS            TLSConfig     `mapstructure:"tls"`
}

type TLSConfig struct {
	Enabled  bool   `mapstructure:"enabled"`
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type MCPConfig struct {
	Enabled   bool              `mapstructure:"enabled"`
	Host      string            `mapstructure:"host"`
	Port      int               `mapstructure:"port"`
	WebSocket WebSocketConfig   `mapstructure:"websocket"`
	CORS      CORSConfig        `mapstructure:"cors"`
}

type WebSocketConfig struct {
	ReadBufferSize  int  `mapstructure:"read_buffer_size"`
	WriteBufferSize int  `mapstructure:"write_buffer_size"`
	CheckOrigin     bool `mapstructure:"check_origin"`
}

type CORSConfig struct {
	AllowedOrigins   []string `mapstructure:"allowed_origins"`
	AllowedMethods   []string `mapstructure:"allowed_methods"`
	AllowedHeaders   []string `mapstructure:"allowed_headers"`
	ExposedHeaders   []string `mapstructure:"exposed_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
	MaxAge           int      `mapstructure:"max_age"`
}

type ProvidersConfig struct {
	Default    string                         `mapstructure:"default"`
	OpenRouter ProviderConfig                 `mapstructure:"openrouter"`
	OpenAI     ProviderConfig                 `mapstructure:"openai"`
	Anthropic  ProviderConfig                 `mapstructure:"anthropic"`
	Ollama     ProviderConfig                 `mapstructure:"ollama"`
}

type ProviderConfig struct {
	APIKey       string        `mapstructure:"api_key"`
	BaseURL      string        `mapstructure:"base_url"`
	DefaultModel string        `mapstructure:"default_model"`
	Timeout      time.Duration `mapstructure:"timeout"`
	RetryCount   int           `mapstructure:"retry_count"`
	RetryDelay   time.Duration `mapstructure:"retry_delay"`
	Models       []string      `mapstructure:"models"`
}

type DatabaseConfig struct {
	Type     string         `mapstructure:"type"`
	SQLite   SQLiteConfig   `mapstructure:"sqlite"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	MySQL    MySQLConfig    `mapstructure:"mysql"`
}

type SQLiteConfig struct {
	Path               string        `mapstructure:"path"`
	MaxConnections     int           `mapstructure:"max_connections"`
	BusyTimeout        time.Duration `mapstructure:"busy_timeout"`
	JournalMode        string        `mapstructure:"journal_mode"`
	Synchronous        string        `mapstructure:"synchronous"`
}

type PostgresConfig struct {
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Database              string        `mapstructure:"database"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	SSLMode               string        `mapstructure:"ssl_mode"`
	MaxConnections        int           `mapstructure:"max_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}

type MySQLConfig struct {
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Database              string        `mapstructure:"database"`
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	MaxConnections        int           `mapstructure:"max_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}

type CacheConfig struct {
	Enabled         bool             `mapstructure:"enabled"`
	Type            string           `mapstructure:"type"`
	TTL             time.Duration    `mapstructure:"ttl"`
	MaxSize         int              `mapstructure:"max_size"`
	CleanupInterval time.Duration    `mapstructure:"cleanup_interval"`
	Redis           RedisConfig      `mapstructure:"redis"`
	Memcached       MemcachedConfig  `mapstructure:"memcached"`
}

type RedisConfig struct {
	Host               string        `mapstructure:"host"`
	Port               int           `mapstructure:"port"`
	Password           string        `mapstructure:"password"`
	DB                 int           `mapstructure:"db"`
	PoolSize           int           `mapstructure:"pool_size"`
	MinIdleConnections int           `mapstructure:"min_idle_connections"`
	DialTimeout        time.Duration `mapstructure:"dial_timeout"`
	ReadTimeout        time.Duration `mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `mapstructure:"write_timeout"`
}

type MemcachedConfig struct {
	Servers              []string      `mapstructure:"servers"`
	Timeout              time.Duration `mapstructure:"timeout"`
	MaxIdleConnections   int           `mapstructure:"max_idle_connections"`
}

type PerformanceConfig struct {
	WorkerPoolSize         int           `mapstructure:"worker_pool_size"`
	BatchSize              int           `mapstructure:"batch_size"`
	BatchDelay             time.Duration `mapstructure:"batch_delay"`
	MaxConcurrentRequests  int           `mapstructure:"max_concurrent_requests"`
	RequestTimeout         time.Duration `mapstructure:"request_timeout"`
	UltraCompressedRatio   float64       `mapstructure:"ultra_compressed_ratio"`
	ThinkingTokens         ThinkingTokensConfig `mapstructure:"thinking_tokens"`
	MaxMemoryMB            int           `mapstructure:"max_memory_mb"`
	MaxGoroutines          int           `mapstructure:"max_goroutines"`
}

type ThinkingTokensConfig struct {
	Standard int `mapstructure:"standard"`
	Deep     int `mapstructure:"deep"`
	Ultra    int `mapstructure:"ultra"`
}

type RateLimitConfig struct {
	Enabled    bool                 `mapstructure:"enabled"`
	Global     RateLimitRule        `mapstructure:"global"`
	PerSession RateLimitRule        `mapstructure:"per_session"`
	PerIP      RateLimitRule        `mapstructure:"per_ip"`
}

type RateLimitRule struct {
	RequestsPerMinute int `mapstructure:"requests_per_minute"`
	Burst             int `mapstructure:"burst"`
}

type SecurityConfig struct {
	APIKeyEncryption  bool       `mapstructure:"api_key_encryption"`
	SessionEncryption bool       `mapstructure:"session_encryption"`
	CORS              CORSConfig `mapstructure:"cors"`
	Auth              AuthConfig `mapstructure:"auth"`
	TLS               TLSSecurityConfig `mapstructure:"tls"`
}

type AuthConfig struct {
	SessionTimeout       time.Duration `mapstructure:"session_timeout"`
	JWTSecret            string        `mapstructure:"jwt_secret"`
	JWTExpiry            time.Duration `mapstructure:"jwt_expiry"`
	RefreshTokenExpiry   time.Duration `mapstructure:"refresh_token_expiry"`
}

type TLSSecurityConfig struct {
	MinVersion    string   `mapstructure:"min_version"`
	CipherSuites  []string `mapstructure:"cipher_suites"`
}

type LoggingConfig struct {
	Level             string                    `mapstructure:"level"`
	Format            string                    `mapstructure:"format"`
	Output            string                    `mapstructure:"output"`
	File              LogFileConfig             `mapstructure:"file"`
	StructuredFields  map[string]string         `mapstructure:"structured_fields"`
	Components        map[string]string         `mapstructure:"components"`
}

type LogFileConfig struct {
	Path        string `mapstructure:"path"`
	MaxSize     string `mapstructure:"max_size"`
	MaxBackups  int    `mapstructure:"max_backups"`
	MaxAge      string `mapstructure:"max_age"`
	Compress    bool   `mapstructure:"compress"`
}

type MonitoringConfig struct {
	Enabled     bool               `mapstructure:"enabled"`
	Metrics     MetricsConfig      `mapstructure:"metrics"`
	Tracing     TracingConfig      `mapstructure:"tracing"`
	HealthCheck HealthCheckConfig  `mapstructure:"health_check"`
	Profiling   ProfilingConfig    `mapstructure:"profiling"`
}

type MetricsConfig struct {
	Enabled   bool   `mapstructure:"enabled"`
	Path      string `mapstructure:"path"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
}

type TracingConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	Provider    string  `mapstructure:"provider"`
	Endpoint    string  `mapstructure:"endpoint"`
	ServiceName string  `mapstructure:"service_name"`
	SampleRate  float64 `mapstructure:"sample_rate"`
}

type HealthCheckConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	Path     string        `mapstructure:"path"`
	Interval time.Duration `mapstructure:"interval"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type ProfilingConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	Path    string `mapstructure:"path"`
}

type SuperClaudeSpecificConfig struct {
	Commands CommandsConfig `mapstructure:"commands"`
	Personas PersonasConfig `mapstructure:"personas"`
	Flags    FlagsConfig    `mapstructure:"flags"`
}

type CommandsConfig struct {
	Enabled              bool                      `mapstructure:"enabled"`
	DefaultPersona       string                    `mapstructure:"default_persona"`
	AutoPersonaSelection bool                      `mapstructure:"auto_persona_selection"`
	CommandHistorySize   int                       `mapstructure:"command_history_size"`
	Analyze              AnalyzeCommandConfig      `mapstructure:"analyze"`
	Build                BuildCommandConfig        `mapstructure:"build"`
	Test                 TestCommandConfig         `mapstructure:"test"`
	Improve              ImproveCommandConfig      `mapstructure:"improve"`
}

type AnalyzeCommandConfig struct {
	MaxFileSize          string   `mapstructure:"max_file_size"`
	SupportedExtensions  []string `mapstructure:"supported_extensions"`
}

type BuildCommandConfig struct {
	Timeout        time.Duration `mapstructure:"timeout"`
	ParallelBuilds bool          `mapstructure:"parallel_builds"`
}

type TestCommandConfig struct {
	Timeout           time.Duration `mapstructure:"timeout"`
	CoverageThreshold int           `mapstructure:"coverage_threshold"`
}

type ImproveCommandConfig struct {
	MaxSuggestions  int  `mapstructure:"max_suggestions"`
	IncludeExamples bool `mapstructure:"include_examples"`
}

type PersonasConfig struct {
	Enabled           bool `mapstructure:"enabled"`
	AllowCustom       bool `mapstructure:"allow_custom"`
	CollaborationMode bool `mapstructure:"collaboration_mode"`
}

type FlagsConfig struct {
	UltraCompressedDefault bool   `mapstructure:"ultra_compressed_default"`
	ThinkingModeDefault    string `mapstructure:"thinking_mode_default"`
	EvidenceModeDefault    bool   `mapstructure:"evidence_mode_default"`
}

type IDEConfig struct {
	Enabled bool           `mapstructure:"enabled"`
	VSCode  VSCodeConfig   `mapstructure:"vscode"`
	Cursor  CursorConfig   `mapstructure:"cursor"`
	Vim     VimConfig      `mapstructure:"vim"`
	Emacs   EmacsConfig    `mapstructure:"emacs"`
}

type VSCodeConfig struct {
	ExtensionID  string `mapstructure:"extension_id"`
	AutoComplete bool   `mapstructure:"auto_complete"`
	CodeActions  bool   `mapstructure:"code_actions"`
	StatusBar    bool   `mapstructure:"status_bar"`
}

type CursorConfig struct {
	Enabled     bool `mapstructure:"enabled"`
	Keybindings bool `mapstructure:"keybindings"`
	ContextMenu bool `mapstructure:"context_menu"`
}

type VimConfig struct {
	PluginName string `mapstructure:"plugin_name"`
	LeaderKey  string `mapstructure:"leader_key"`
}

type EmacsConfig struct {
	PackageName string `mapstructure:"package_name"`
	PrefixKey   string `mapstructure:"prefix_key"`
}

type FeaturesConfig struct {
	MCPServer          bool `mapstructure:"mcp_server"`
	CacheOptimization  bool `mapstructure:"cache_optimization"`
	BatchProcessing    bool `mapstructure:"batch_processing"`
	ParallelExecution  bool `mapstructure:"parallel_execution"`
	CommandCompletion  bool `mapstructure:"command_completion"`
	SessionPersistence bool `mapstructure:"session_persistence"`
	MetricsCollection  bool `mapstructure:"metrics_collection"`
	AutoUpdates        bool `mapstructure:"auto_updates"`
}

type DevelopmentConfig struct {
	Debug        bool             `mapstructure:"debug"`
	HotReload    bool             `mapstructure:"hot_reload"`
	MockProviders bool            `mapstructure:"mock_providers"`
	TestMode     bool             `mapstructure:"test_mode"`
	Profiling    bool             `mapstructure:"profiling"`
	Fixtures     FixturesConfig   `mapstructure:"fixtures"`
}

type FixturesConfig struct {
	LoadTestData bool   `mapstructure:"load_test_data"`
	TestDataPath string `mapstructure:"test_data_path"`
}

type DeploymentConfig struct {
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`
	BuildTime   string `mapstructure:"build_time"`
	GitCommit   string `mapstructure:"git_commit"`
}

// LoadConfig loads configuration from files and environment variables
func LoadConfig(configPath string) (*SuperClaudeConfig, error) {
	v := viper.New()
	
	// Set defaults
	setAdvancedDefaults(v)
	
	// Set config file name and paths
	v.SetConfigName("superclaude")
	v.SetConfigType("yaml")
	
	// Add config paths
	if configPath != "" {
		v.AddConfigPath(configPath)
	}
	v.AddConfigPath("./config")
	v.AddConfigPath("$HOME/.superclaude")
	v.AddConfigPath("/etc/superclaude")
	
	// Environment variable configuration
	v.SetEnvPrefix("SUPERCLAUDE")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// Read main config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}
	
	// Read environment-specific config
	environment := v.GetString("deployment.environment")
	if environment == "" {
		environment = os.Getenv("ENVIRONMENT")
		if environment == "" {
			environment = "development"
		}
	}
	
	// Merge environment-specific config
	if err := mergeEnvironmentConfig(v, environment); err != nil {
		return nil, fmt.Errorf("error merging environment config: %w", err)
	}
	
	// Unmarshal to struct
	var config SuperClaudeConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	return &config, nil
}

// mergeEnvironmentConfig merges environment-specific configuration
func mergeEnvironmentConfig(v *viper.Viper, environment string) error {
	envConfigFile := filepath.Join(v.ConfigFileUsed(), "..", environment+".yaml")
	if _, err := os.Stat(envConfigFile); err == nil {
		envViper := viper.New()
		envViper.SetConfigFile(envConfigFile)
		if err := envViper.ReadInConfig(); err != nil {
			return err
		}
		
		// Merge environment config
		for key, value := range envViper.AllSettings() {
			v.Set(key, value)
		}
	}
	
	return nil
}

// setAdvancedDefaults sets default configuration values
func setAdvancedDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.timeout", "30s")
	v.SetDefault("server.max_connections", 1000)
	
	// Database defaults
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.sqlite.path", "~/.superclaude/superclaude.db")
	
	// Cache defaults
	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.type", "memory")
	v.SetDefault("cache.ttl", "15m")
	v.SetDefault("cache.max_size", 1000)
	
	// Performance defaults
	v.SetDefault("performance.worker_pool_size", 0)
	v.SetDefault("performance.batch_size", 10)
	v.SetDefault("performance.batch_delay", "100ms")
	
	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	
	// Features defaults
	v.SetDefault("features.mcp_server", true)
	v.SetDefault("features.cache_optimization", true)
	v.SetDefault("features.batch_processing", true)
}

// validateConfig validates the configuration
func validateConfig(config *SuperClaudeConfig) error {
	// Validate required fields
	if config.Providers.Default == "" {
		return fmt.Errorf("providers.default is required")
	}
	
	// Validate port ranges
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}
	
	// Validate database configuration
	switch config.Database.Type {
	case "sqlite", "postgres", "mysql":
		// Valid
	default:
		return fmt.Errorf("database.type must be one of: sqlite, postgres, mysql")
	}
	
	// Validate cache configuration
	switch config.Cache.Type {
	case "memory", "redis", "memcached":
		// Valid
	default:
		return fmt.Errorf("cache.type must be one of: memory, redis, memcached")
	}
	
	return nil
}