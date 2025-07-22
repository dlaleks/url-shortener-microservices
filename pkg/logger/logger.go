package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger wraps zap.Logger with additional functionality
type Logger struct {
	*zap.Logger
}

// Config holds logger configuration
type Config struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"`         // json, text
	Output     string `json:"output" yaml:"output"`         // stdout, file, both
	FilePath   string `json:"file_path" yaml:"file_path"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`       // megabytes
	MaxBackups int    `json:"max_backups" yaml:"max_backups"` // number of old files
	MaxAge     int    `json:"max_age" yaml:"max_age"`         // days
	Compress   bool   `json:"compress" yaml:"compress"`
	ServiceName string `json:"service_name" yaml:"service_name"`
}

// NewLogger creates a new logger instance
func NewLogger(config Config) (*Logger, error) {
	// Parse log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	// Create encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "message"
	encoderConfig.StacktraceKey = "stacktrace"

	// Create encoder
	var encoder zapcore.Encoder
	switch config.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "text", "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create writers
	var cores []zapcore.Core

	switch config.Output {
	case "stdout":
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	case "file":
		if config.FilePath == "" {
			return nil, fmt.Errorf("file_path is required when output is 'file'")
		}
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level))
	case "both":
		// Stdout
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
		// File
		if config.FilePath != "" {
			fileWriter := &lumberjack.Logger{
				Filename:   config.FilePath,
				MaxSize:    config.MaxSize,
				MaxBackups: config.MaxBackups,
				MaxAge:     config.MaxAge,
				Compress:   config.Compress,
			}
			cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileWriter), level))
		}
	default:
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// Combine cores
	core := zapcore.NewTee(cores...)

	// Create logger with options
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	// Add service name field if provided
	if config.ServiceName != "" {
		options = append(options, zap.Fields(zap.String("service", config.ServiceName)))
	}

	logger := zap.New(core, options...)

	return &Logger{Logger: logger}, nil
}

// Default creates a logger with default configuration
func Default(serviceName string) *Logger {
	config := Config{
		Level:       "info",
		Format:      "json",
		Output:      "stdout",
		ServiceName: serviceName,
	}

	logger, err := NewLogger(config)
	if err != nil {
		// Fallback to basic logger
		zapLogger, _ := zap.NewProduction()
		return &Logger{Logger: zapLogger}
	}

	return logger
}

// WithError adds error field to log entry
func (l *Logger) WithError(err error) *zap.Logger {
	return l.With(zap.Error(err))
}

// WithField adds a field to log entry
func (l *Logger) WithField(key string, value interface{}) *zap.Logger {
	return l.With(zap.Any(key, value))
}

// WithFields adds multiple fields to log entry
func (l *Logger) WithFields(fields map[string]interface{}) *zap.Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return l.With(zapFields...)
}

// WithRequestID adds request ID to log entry
func (l *Logger) WithRequestID(requestID string) *zap.Logger {
	return l.With(zap.String("request_id", requestID))
}

// WithUserID adds user ID to log entry
func (l *Logger) WithUserID(userID string) *zap.Logger {
	return l.With(zap.String("user_id", userID))
}

// WithComponent adds component name to log entry
func (l *Logger) WithComponent(component string) *zap.Logger {
	return l.With(zap.String("component", component))
}

// HTTP middleware logging helper
func (l *Logger) HTTPMiddleware() *zap.Logger {
	return l.Named("http")
}

// Database logging helper
func (l *Logger) Database() *zap.Logger {
	return l.Named("database")
}

// gRPC logging helper
func (l *Logger) GRPC() *zap.Logger {
	return l.Named("grpc")
}

// Cache logging helper
func (l *Logger) Cache() *zap.Logger {
	return l.Named("cache")
}

// Queue logging helper
func (l *Logger) Queue() *zap.Logger {
	return l.Named("queue")
}

// External service logging helper
func (l *Logger) External(serviceName string) *zap.Logger {
	return l.Named("external").With(zap.String("external_service", serviceName))
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.Logger.Sync()
}

// Close closes the logger
func (l *Logger) Close() error {
	return l.Sync()
}

// Global logger instance
var globalLogger *Logger

// SetGlobalLogger sets the global logger instance
func SetGlobalLogger(logger *Logger) {
	globalLogger = logger
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		globalLogger = Default("app")
	}
	return globalLogger
}

// Global convenience methods
func Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	GetGlobalLogger().Panic(msg, fields...)
}