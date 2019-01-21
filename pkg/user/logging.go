package user

import (
	"context"
	"log"
	"time"
)

type loggingService struct {
	logger *log.Logger
	next   Service
}

// WrapWithLogging wraps an WorkspaceService and logs the outputs.
func WrapWithLogging(logger *log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Register(ctx context.Context, input RegisterUserInput) (jwt JWT, err error) {
	defer func(begin time.Time) {
		s.logger.Printf("method: Register, input: %+v, took: %s, error: %v", input, time.Since(begin), err)
	}(time.Now())
	return s.next.Register(ctx, input)
}
