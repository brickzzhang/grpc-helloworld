package logger

import (
	"context"

	"go.uber.org/zap"
)

// Info info level log
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Info(msg, fields...)
}

// Warn warn level log
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Warn(msg, fields...)
}

// Error error level log
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Error(msg, fields...)
}
