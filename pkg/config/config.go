package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Log configuration
	Log LogConfig `mapstructure:"log"`

	// Server configuration
	Server ServerConfig `mapstructure:"server"`

	// Database configuration
	Database DatabaseConfig `mapstructure:"database"`

	// NLP configuration
	NLP NLPConfig `mapstructure:"nlp"`

	// Telemetry configuration
	Telemetry TelemetryConfig `mapstructure:"telemetry"`

	// Embedding configuration
	Embedding EmbeddingConfig `mapstructure:"embedding"`

	// Alert configuration
	Alert AlertConfig `mapstructure:"alert"`

	// CircuitBreaker configuration
	CircuitBreaker CircuitBreakerConfig `mapstructure:"circuit_breaker"`
}

// AlertConfig holds configuration for alerting
type AlertConfig struct {
	Enabled  bool     `mapstructure:"enabled"`
	SMTPHost string   `mapstructure:"smtp_host"`
	SMTPPort int      `mapstructure:"smtp_port"`
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	From     string   `mapstructure:"from"`
	To       []string `mapstructure:"to"`
}

// CircuitBreakerConfig holds configuration for circuit breaking
type CircuitBreakerConfig struct {
	Enabled          bool    `mapstructure:"enabled"`
	MaxRequests      uint32  `mapstructure:"max_requests"`
	Interval         int     `mapstructure:"interval"` // in seconds
	Timeout          int     `mapstructure:"timeout"`  // in seconds
	ReadyToTripRatio float64 `mapstructure:"ready_to_trip_ratio"`
}

// TelemetryConfig holds telemetry configuration
type TelemetryConfig struct {
	ParquetPath string `mapstructure:"parquet_path"`
	DbURL       string `mapstructure:"db_url"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"` // gin mode: debug, release, test
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string `mapstructure:"driver"` // neo4j, falkordb
	URI      string `mapstructure:"uri"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
}

// NLPConfig holds NLP configuration
type NLPConfig struct {
	// Models is a map of model configurations (e.g. "default", "embedding", "summary")
	Models map[string]NLPModelConfig `mapstructure:"models"`

	// RouterRules defines how to route requests
	RouterRules []RouterRule `mapstructure:"router_rules"`
}

// NLPModelConfig holds configuration for a specific model
type NLPModelConfig struct {
	Provider    string  `mapstructure:"provider"` // type: openai, anthropic
	Model       string  `mapstructure:"model"`
	APIKey      string  `mapstructure:"api_key"`
	BaseURL     string  `mapstructure:"base_url"`
	Temperature float32 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
}

// RouterRule defines a rule for routing requests
type RouterRule struct {
	Usage    string `mapstructure:"usage"`    // Tag to match (e.g. "hipaa", "coding")
	Provider string `mapstructure:"provider"` // Provider ID to use
	Fallback string `mapstructure:"fallback"` // Fallback provider ID
}

// EmbeddingConfig holds embedding configuration
type EmbeddingConfig struct {
	Provider string `mapstructure:"provider"` // openai, etc.
	Model    string `mapstructure:"model"`
	APIKey   string `mapstructure:"api_key"`
	BaseURL  string `mapstructure:"base_url"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	// Set defaults
	setDefaults()

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Override with environment variables if present
	overrideWithEnv(config)

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")

	// Server defaults
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.driver", "ladybug")
	viper.SetDefault("database.uri", "./ladybug_db")
	viper.SetDefault("database.username", "")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "")

	viper.SetDefault("nlp.models.default.provider", "rustbert")
	viper.SetDefault("nlp.models.default.base_url", "rustbert://generator")
	viper.SetDefault("nlp.models.default.model", "gpt2")
	viper.SetDefault("nlp.models.default.temperature", 0.7)
	viper.SetDefault("nlp.models.default.max_tokens", 256)

	viper.SetDefault("nlp.models.embedding.provider", "embedeverything")
	viper.SetDefault("nlp.models.embedding.base_url", "embedeverything://")
	viper.SetDefault("nlp.models.embedding.model", "all-MiniLM-L6-v2")

	// Telemetry defaults
	home, err := os.UserHomeDir()
	if err == nil {
		defaultPath := fmt.Sprintf("%s/.predicato/telemetry", home)
		viper.SetDefault("telemetry.parquet_path", defaultPath)
	}
}

// overrideWithEnv overrides config with environment variables
func overrideWithEnv(config *Config) {
	// Initialize Models map if nil
	if config.NLP.Models == nil {
		config.NLP.Models = make(map[string]NLPModelConfig)
	}

	// Helper to get or create model config
	getModel := func(name string) NLPModelConfig {
		if c, ok := config.NLP.Models[name]; ok {
			return c
		}
		return NLPModelConfig{}
	}

	// Update default model from env
	defaultModel := getModel("default")
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		defaultModel.APIKey = apiKey
	}
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" && defaultModel.Provider == "anthropic" {
		defaultModel.APIKey = apiKey
	}
	config.NLP.Models["default"] = defaultModel

	// Update embedding model from env
	embeddingModel := getModel("embedding")
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		embeddingModel.APIKey = apiKey
	}
	config.NLP.Models["embedding"] = embeddingModel

	// Database credentials
	if uri := os.Getenv("NEO4J_URI"); uri != "" {
		config.Database.URI = uri
	}
	if user := os.Getenv("NEO4J_USER"); user != "" {
		config.Database.Username = user
	}
	if pass := os.Getenv("NEO4J_PASSWORD"); pass != "" {
		config.Database.Password = pass
	}

	// ladybug database path
	if dbPath := os.Getenv("ladybug_DB_PATH"); dbPath != "" {
		config.Database.URI = dbPath
	}

	// Generic database settings
	if dbDriver := os.Getenv("DB_DRIVER"); dbDriver != "" {
		config.Database.Driver = dbDriver
	}
	if dbURI := os.Getenv("DB_URI"); dbURI != "" {
		config.Database.URI = dbURI
	}

	// Server settings
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if port := os.Getenv("SERVER_PORT"); port != "" {
		viper.Set("server.port", port)
	}

	// Telemetry settings
	if path := os.Getenv("TELEMETRY_PARQUET_PATH"); path != "" {
		config.Telemetry.ParquetPath = path
	}
}
