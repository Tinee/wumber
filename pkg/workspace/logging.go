package workspace

import (
	"context"
	"log"
	"time"
	"wumber"
)

type loggingService struct {
	logger *log.Logger
	next   Service
}

// WrapWithLogging wraps an WorkspaceService and logs the outputs.
func WrapWithLogging(logger *log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, name, accountID string) (id wumber.WorkspaceID, err error) {
	defer func(begin time.Time) {
		s.logger.Printf("method: create, name: %s, took: %s, error: %s", name, time.Since(begin), err)
	}(time.Now())
	return s.next.Create(ctx, name, accountID)
}
