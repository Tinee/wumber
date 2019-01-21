package workspace

import (
	"context"
	"wumber"
)

type Service interface {
	Create(ctx context.Context, name, accountID string) (wumber.WorkspaceID, error)
}

type service struct {
	wsRepo wumber.WorkspaceRepository
}

func NewService(r wumber.WorkspaceRepository) Service {
	return &service{
		wsRepo: r,
	}
}

func (s *service) Create(ctx context.Context, name, accountID string) (wumber.WorkspaceID, error) {
	if name == "" || accountID == "" {
		return "", ErrInvalidArgument
	}
	id, err := s.wsRepo.CreateWorkspace(ctx, name, accountID)
	if err != nil {
		return "", err
	}

	return id, nil
}
