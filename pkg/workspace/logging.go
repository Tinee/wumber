package workspace

import (
	"context"
	"time"
	"wumber"
	"wumber/logger"

	logrus "github.com/sirupsen/logrus"
)

type loggingService struct {
	logger *logger.Logger
	next   Service
}

// WrapWithLogging wraps an WorkspaceService and logs the outputs.
func WrapWithLogging(logger *logger.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Create(ctx context.Context, input CreateWorkspaceInput) (id wumber.WorkspaceID, err error) {
	s.logger.WithFields(logrus.Fields{
		"input": input,
	}).Debug("Calling Create")

	defer func(begin time.Time) {
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"input": input,
				"took":  time.Since(begin),
			}).WithError(err).Error("Error when calling Create")
			return
		}

		s.logger.WithFields(logrus.Fields{
			"input":  input,
			"output": id,
			"took":   time.Since(begin),
		}).Debug("Successfully created the Workspace.")
	}(time.Now())
	return s.next.Create(ctx, input)
}
