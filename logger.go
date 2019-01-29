package wumber

import "context"

// Logger is the domain for when we want to log in the application.
type Logger interface {
	Error(ctx context.Context, msg string, args ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
	Flush() error
}
