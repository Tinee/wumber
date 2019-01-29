package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
	contextMappers []ContextValueMapperFunc
}

// Debug wraps a third party logger, it adds metadata from the context into the error logs.
func (l *Logger) Error(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, l.applyContextMappers(ctx)...)
	l.With(args...).Error(msg)
}

// Debug wraps a third party logger, it adds metadata from the context into the debug logs.
func (l *Logger) Debug(ctx context.Context, msg string, args ...interface{}) {
	args = append(args, l.applyContextMappers(ctx)...)
	l.With(args...).Debug(msg)
}

// Flush flushes any buffered log entries.
func (l *Logger) Flush() error { return l.Sync() }

// NewLogger determines by looking at the environment, which level the logger should log with.
// If the environment isn't set it will default to a Developer logger.
func NewLogger(env string, funcs ...ContextValueMapperFunc) *Logger {
	var logger *zap.Logger
	switch strings.ToLower(env) {
	case "uat", "dev":
		logger, _ = zap.NewDevelopment()
	case "prod":
		logger, _ = zap.NewProduction()
	default:
		logger, _ = zap.NewProduction()
	}

	return &Logger{SugaredLogger: logger.Sugar(), contextMappers: funcs}
}

// NewAWSLogger is a shortcut for logger.NewLogger("dev", logger.AWSRequestID)
func NewAWSLogger(env string) *Logger {
	return NewLogger(env, AWSRequestID)
}
