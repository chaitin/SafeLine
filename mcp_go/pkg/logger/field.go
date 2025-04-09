package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Field 日志字段
type Field = zapcore.Field

// String 创建字符串字段
func String(key string, val string) Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) Field {
	return zap.Int(key, val)
}

// Int64 创建 int64 字段
func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

// Float64 创建浮点数字段
func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

// Bool 创建布尔字段
func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

// Err 创建错误字段
func Err(err error) Field {
	return zap.Error(err)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) Field {
	return zap.Any(key, val)
}

// Duration 创建时间段字段
func Duration(key string, val float64) Field {
	return zap.Float64(key, val)
}
