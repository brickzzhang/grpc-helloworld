// Package logger provides log package
package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logType string

// Logger interface for different levels of log
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

func defaultZapConfig() zap.Config {
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

// Get default zap logger with json format
func getDefaultZapLogger() *zap.Logger {
	zapLogger, err := defaultZapConfig().Build()
	if err != nil {
		return zap.NewExample()
	}
	return zapLogger
}

// NewLoggerToCtx insert logger to context
func NewLoggerToCtx(ctx context.Context, log Logger) context.Context {
	if log == nil {
		return context.WithValue(ctx, logType("logger"), getDefaultZapLogger())
	}

	return context.WithValue(ctx, logType("logger"), log)
}

func extractLoggerFromCtx(ctx context.Context) Logger {
	v, ok := ctx.Value(logType("logger")).(Logger)
	if !ok {
		return nil
	}

	return v
}
