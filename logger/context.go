package logger

import (
	"context"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// ContextValueMapperFunc takes a context. if it finds the value it return the key & value.
// If we don't find anything we return ok false.
type ContextValueMapperFunc func(ctx context.Context) (key string, value string, ok bool)

// AWSRequestID implements the ContextValueMapperFunc interface, it finds requestId in the context that AWS provide.
func AWSRequestID(ctx context.Context) (string, string, bool) {
	context, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return "", "", false
	}

	return "requestId", context.AwsRequestID, true
}

func (l *Logger) applyContextMappers(ctx context.Context) []interface{} {
	var result []interface{}
	for _, fn := range l.contextMappers {
		key, value, ok := fn(ctx)
		if !ok {
			continue
		}
		result = append(result, key, value)
	}
	return result
}
