package user

import (
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
