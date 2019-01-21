package user

import (
	"context"
	"time"
	"wumber/logger"

	log "github.com/sirupsen/logrus"
)

type loggingService struct {
	logger *logger.Logger
	next   Service
}

// WrapWithLogging wraps an WorkspaceService and logs the outputs.
func WrapWithLogging(logger *logger.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Register(ctx context.Context, input RegisterUserInput) (jwt JWT, err error) {
	s.logger.WithFields(log.Fields{
		"input": input,
	}).Debug("Calling Register.")

	defer func(begin time.Time) {
		if err != nil {
			s.logger.WithFields(log.Fields{
				"took": time.Since(begin),
			}).WithError(err).Error("Called Register and got an error.")
			return
		}

		s.logger.WithFields(log.Fields{
			"input":  input,
			"output": jwt,
			"took":   time.Since(begin),
		}).Debug("Called Register.")
	}(time.Now())

	return s.next.Register(ctx, input)
}
