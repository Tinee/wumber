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

func (s *loggingService) Create(ctx context.Context, name, accountID string) (id wumber.WorkspaceID, err error) {
	s.logger.WithFields(logrus.Fields{
		"name":      name,
		"accountId": accountID,
	}).Debug("Calling Create")

	defer func(begin time.Time) {
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"name":      name,
				"accountID": name,
				"took":      time.Since(begin),
			}).WithError(err).Error("Error when calling Create")
			return
		}

		s.logger.WithFields(logrus.Fields{
			"name":      name,
			"accountId": accountID,
			"output":    id,
			"took":      time.Since(begin),
		}).Debug("Called Create")
	}(time.Now())
	return s.next.Create(ctx, name, accountID)
}
