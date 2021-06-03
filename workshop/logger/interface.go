package logger

import (
	"context"

	"go.uber.org/zap"
)

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	extractLoggerFromCtx(ctx).Error(msg, fields...)
}
