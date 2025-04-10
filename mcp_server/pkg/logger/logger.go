package logger

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 封装 zap.Logger
type Logger struct {
	zl *zap.Logger
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// Config 日志配置
type Config struct {
	// 日志级别
	Level string `json:"level" yaml:"level"`
	// 日志文件路径
	FilePath string `json:"file_path" yaml:"file_path"`
	// 是否输出到控制台
	Console bool `json:"console" yaml:"console"`
	// 是否记录调用者信息
	Caller bool `json:"caller" yaml:"caller"`
	// 是否使用开发模式（更详细的日志）
	Development bool `json:"development" yaml:"development"`
}

// 默认配置
var defaultConfig = Config{
	Level:       "info",
	FilePath:    "logs/mcp.log",
	Console:     true,
	Caller:      true,
	Development: false,
}

// Init 初始化日志系统
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		if cfg == nil {
			cfg = &defaultConfig
		}

		// 确保日志目录存在
		if cfg.FilePath != "" {
			dir := filepath.Dir(cfg.FilePath)
			if err = os.MkdirAll(dir, 0755); err != nil {
				return
			}
		}

		// 配置编码器
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stacktrace",
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
			},
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// 设置日志级别
		var level zapcore.Level
		if err = level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}

		// 创建Core
		var cores []zapcore.Core

		// 文件输出
		if cfg.FilePath != "" {
			fileWriter, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			cores = append(cores, zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(fileWriter),
				level,
			))
		}

		// 控制台输出
		if cfg.Console {
			cores = append(cores, zapcore.NewCore(
				zapcore.NewConsoleEncoder(encoderConfig),
				zapcore.AddSync(os.Stdout),
				level,
			))
		}

		// 创建Logger
		core := zapcore.NewTee(cores...)
		zl := zap.New(core)

		// 是否记录调用者信息
		if cfg.Caller {
			zl = zl.WithOptions(zap.AddCaller())
		}

		// 是否使用开发模式
		if cfg.Development {
			zl = zl.WithOptions(zap.Development())
		}

		defaultLogger = &Logger{zl: zl}
	})

	return err
}

// GetLogger 获取日志实例
func GetLogger() *Logger {
	if defaultLogger == nil {
		Init(nil)
	}
	return defaultLogger
}

// With 创建带有字段的新Logger
func (l *Logger) With(key string, value interface{}) *Logger {
	return &Logger{zl: l.zl.With(Any(key, value))}
}

// Debug level
func (l *Logger) Debug(msg string, fields ...Field) {
	l.zl.Debug(msg, fields...)
}

// Info level
func (l *Logger) Info(msg string, fields ...Field) {
	l.zl.Info(msg, fields...)
}

// Warn level
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zl.Warn(msg, fields...)
}

// Error level
func (l *Logger) Error(msg string, fields ...Field) {
	l.zl.Error(msg, fields...)
}

// Fatal level
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zl.Fatal(msg, fields...)
}

// 全局函数

// With 创建带有字段的新Logger
func With(key string, value interface{}) *Logger {
	l := GetLogger()
	// 为链式调用创建新的logger实例，并添加caller skip
	return &Logger{zl: l.zl.WithOptions(zap.AddCallerSkip(1)).With(Any(key, value))}
}

// Debug level
func Debug(msg string, fields ...Field) {
	l := GetLogger()
	l.zl.WithOptions(zap.AddCallerSkip(1)).Debug(msg, fields...)
}

// Info level
func Info(msg string, fields ...Field) {
	l := GetLogger()
	l.zl.WithOptions(zap.AddCallerSkip(1)).Info(msg, fields...)
}

// Warn level
func Warn(msg string, fields ...Field) {
	l := GetLogger()
	l.zl.WithOptions(zap.AddCallerSkip(1)).Warn(msg, fields...)
}

// Error level
func Error(msg string, fields ...Field) {
	l := GetLogger()
	l.zl.WithOptions(zap.AddCallerSkip(1)).Error(msg, fields...)
}

// Fatal level
func Fatal(msg string, fields ...Field) {
	l := GetLogger()
	l.zl.WithOptions(zap.AddCallerSkip(1)).Fatal(msg, fields...)
}
