package user

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

func (s *loggingService) Register(ctx context.Context, input RegisterUserInput) (jwt wumber.JWT, err error) {
	defer s.logger.Flush()
	s.logger.Debug(ctx, "Calling Register",
		"input", input,
	)
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Error(ctx, "Failed to call Register.",
				"input", input,
				"err", err,
				"took", time.Since(begin),
			)
			return
		}

		s.logger.Debug(ctx, "Called Register",
			"input", input,
			"output", jwt,
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.next.Register(ctx, input)
}
