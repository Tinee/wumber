package workspace

import (
	"context"
	"time"
	"wumber"
)

type loggingService struct {
	logger wumber.Logger
	next   Service
}

// WrapWithLogging wraps an WorkspaceService and logs the outputs.
func WrapWithLogging(logger wumber.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, input CreateWorkspaceInput) (id wumber.WorkspaceID, err error) {
	defer s.logger.Flush()
	s.logger.Debug(ctx, "Calling Create",
		"input", input,
	)
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Error(ctx, "Error when calling Create",
				"input", input,
				"took", time.Since(begin),
			)
			return
		}

		s.logger.Debug(ctx, "Successfully created the Workspace",
			"input", input,
			"output", id,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Create(ctx, input)
}
