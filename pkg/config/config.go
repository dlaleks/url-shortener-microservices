package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	Database     string `mapstructure:"database"`
	SSLMode      string `mapstructure:"ssl_mode"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxLifetime  string `mapstructure:"max_lifetime"`
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	MaxRetries   int    `mapstructure:"max_retries"`
	DialTimeout  string `mapstructure:"dial_timeout"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
	PoolSize     int    `mapstructure:"pool_size"`
}

// MongoConfig holds MongoDB connection configuration
type MongoConfig struct {
	URI             string `mapstructure:"uri"`
	Database        string `mapstructure:"database"`
	ConnectTimeout  string `mapstructure:"connect_timeout"`
	MaxPoolSize     int    `mapstructure:"max_pool_size"`
	MinPoolSize     int    `mapstructure:"min_pool_size"`
	MaxConnIdleTime string `mapstructure:"max_conn_idle_time"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Host              string `mapstructure:"host"`
	Port              int    `mapstructure:"port"`
	ReadTimeout       string `mapstructure:"read_timeout"`
	WriteTimeout      string `mapstructure:"write_timeout"`
	IdleTimeout       string `mapstructure:"idle_timeout"`
	ShutdownTimeout   string `mapstructure:"shutdown_timeout"`
	EnableCORS        bool   `mapstructure:"enable_cors"`
	EnableCompression bool   `mapstructure:"enable_compression"`
	TrustedProxies    []string `mapstructure:"trusted_proxies"`
}

// GRPCConfig holds gRPC server configuration
type GRPCConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	EnableTLS   bool   `mapstructure:"enable_tls"`
	CertFile    string `mapstructure:"cert_file"`
	KeyFile     string `mapstructure:"key_file"`
	EnableAuth  bool   `mapstructure:"enable_auth"`
	AuthSecret  string `mapstructure:"auth_secret"`
}

// JWTConfig holds JWT token configuration
type JWTConfig struct {
	AccessTokenSecret   string `mapstructure:"access_token_secret"`
	RefreshTokenSecret  string `mapstructure:"refresh_token_secret"`
	AccessTokenExpiry   string `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry  string `mapstructure:"refresh_token_expiry"`
	Issuer              string `mapstructure:"issuer"`
	Audience            string `mapstructure:"audience"`
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, text
	Output     string `mapstructure:"output"` // stdout, file
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`    // megabytes
	MaxBackups int    `mapstructure:"max_backups"` // number of old files
	MaxAge     int    `mapstructure:"max_age"`     // days
	Compress   bool   `mapstructure:"compress"`
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Enabled           bool              `mapstructure:"enabled"`
	WindowSize        string            `mapstructure:"window_size"`
	RequestsPerWindow int               `mapstructure:"requests_per_window"`
	SkipSuccessful    bool              `mapstructure:"skip_successful"`
	KeyGenerator      string            `mapstructure:"key_generator"` // ip, user_id, api_key
	Store             string            `mapstructure:"store"`         // memory, redis
	Endpoints         map[string]int    `mapstructure:"endpoints"`     // endpoint-specific limits
}

// MetricsConfig holds metrics/monitoring configuration
type MetricsConfig struct {
	Enabled    bool   `mapstructure:"enabled"`
	Path       string `mapstructure:"path"`
	Namespace  string `mapstructure:"namespace"`
	Subsystem  string `mapstructure:"subsystem"`
	EnablePush bool   `mapstructure:"enable_push"`
	Gateway    string `mapstructure:"gateway"`
	JobLabel   string `mapstructure:"job_label"`
}

// TracingConfig holds distributed tracing configuration
type TracingConfig struct {
	Enabled     bool    `mapstructure:"enabled"`
	ServiceName string  `mapstructure:"service_name"`
	Endpoint    string  `mapstructure:"endpoint"`
	SampleRate  float64 `mapstructure:"sample_rate"`
}

// BaseConfig holds common configuration for all services
type BaseConfig struct {
	Environment string          `mapstructure:"environment"`
	Debug       bool            `mapstructure:"debug"`
	ServiceName string          `mapstructure:"service_name"`
	Version     string          `mapstructure:"version"`
	Log         LogConfig       `mapstructure:"log"`
	Server      ServerConfig    `mapstructure:"server"`
	GRPC        GRPCConfig      `mapstructure:"grpc"`
	Database    DatabaseConfig  `mapstructure:"database"`
	Redis       RedisConfig     `mapstructure:"redis"`
	MongoDB     MongoConfig     `mapstructure:"mongodb"`
	JWT         JWTConfig       `mapstructure:"jwt"`
	RateLimit   RateLimitConfig `mapstructure:"rate_limit"`
	Metrics     MetricsConfig   `mapstructure:"metrics"`
	Tracing     TracingConfig   `mapstructure:"tracing"`
}

// Config interface for service-specific configurations
type Config interface {
	Validate() error
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig[T Config](configPath, serviceName string) (T, error) {
	var config T

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	
	// Add config paths
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/" + serviceName)

	// Environment variables
	viper.SetEnvPrefix(strings.ToUpper(serviceName))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return config, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, continue with defaults and env vars
	}

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.idle_timeout", "120s")
	viper.SetDefault("server.shutdown_timeout", "30s")
	viper.SetDefault("server.enable_cors", true)
	viper.SetDefault("server.enable_compression", true)

	// gRPC defaults
	viper.SetDefault("grpc.host", "0.0.0.0")
	viper.SetDefault("grpc.port", 9090)
	viper.SetDefault("grpc.enable_tls", false)
	viper.SetDefault("grpc.enable_auth", true)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.ssl_mode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.max_lifetime", "5m")

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.max_retries", 3)
	viper.SetDefault("redis.dial_timeout", "5s")
	viper.SetDefault("redis.read_timeout", "3s")
	viper.SetDefault("redis.write_timeout", "3s")
	viper.SetDefault("redis.pool_size", 10)

	// MongoDB defaults
	viper.SetDefault("mongodb.connect_timeout", "10s")
	viper.SetDefault("mongodb.max_pool_size", 100)
	viper.SetDefault("mongodb.min_pool_size", 10)
	viper.SetDefault("mongodb.max_conn_idle_time", "30s")

	// JWT defaults
	viper.SetDefault("jwt.access_token_expiry", "15m")
	viper.SetDefault("jwt.refresh_token_expiry", "168h") // 7 days
	viper.SetDefault("jwt.issuer", "url-shortener")
	viper.SetDefault("jwt.audience", "url-shortener-users")

	// Logging defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 28)
	viper.SetDefault("log.compress", true)

	// Rate limiting defaults
	viper.SetDefault("rate_limit.enabled", true)
	viper.SetDefault("rate_limit.window_size", "1m")
	viper.SetDefault("rate_limit.requests_per_window", 100)
	viper.SetDefault("rate_limit.skip_successful", false)
	viper.SetDefault("rate_limit.key_generator", "ip")
	viper.SetDefault("rate_limit.store", "redis")

	// Metrics defaults
	viper.SetDefault("metrics.enabled", true)
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.namespace", "url_shortener")

	// Tracing defaults
	viper.SetDefault("tracing.enabled", false)
	viper.SetDefault("tracing.sample_rate", 0.1)

	// Environment defaults
	viper.SetDefault("environment", "development")
	viper.SetDefault("debug", false)
	viper.SetDefault("version", "dev")
}

// GetDSN returns database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode)
}

// GetRedisAddr returns Redis connection address
func (r *RedisConfig) GetRedisAddr() string {
	return r.Addr
}

// GetServerAddr returns HTTP server address
func (s *ServerConfig) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// GetGRPCAddr returns gRPC server address
func (g *GRPCConfig) GetGRPCAddr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}