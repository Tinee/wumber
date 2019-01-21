package workspace

import (
	"context"
	"strings"
	"wumber"

	"github.com/pkg/errors"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// Service is the interface that defines the Workspace Service.
type Service interface {
	// Create takes a name and an accountID and attempts to insert it into the database.
	Create(context.Context, CreateWorkspaceInput) (wumber.WorkspaceID, error)
}

type service struct {
	wsRepo wumber.WorkspaceRepository
}

// NewService assembles a Service.
func NewService(r wumber.WorkspaceRepository) Service {
	return &service{
		wsRepo: r,
	}
}

type CreateWorkspaceInput struct {
	OwnerID       string `json:"ownerId"`
	WorkspaceName string `json:"workspaceName"`
}

func (s *service) Create(ctx context.Context, input CreateWorkspaceInput) (wumber.WorkspaceID, error) {
	workspace, err := input.toWorkspace()
	if err != nil {
		return "", errors.Wrap(err, "couldn't convert the input to a Workspace")
	}
	id, err := s.wsRepo.CreateWorkspace(ctx, workspace)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (w CreateWorkspaceInput) toWorkspace() (wumber.Workspace, error) {
	err := validation.ValidateStruct(&w,
		validation.Field(&w.WorkspaceName, validation.Required),
		validation.Field(&w.OwnerID, validation.Required, is.UUID),
	)

	if err != nil {
		return wumber.Workspace{}, err
	}

	return wumber.Workspace{
		Name:    strings.ToTitle(strings.TrimSpace(w.WorkspaceName)),
		OwnerID: w.OwnerID,
	}, nil
}
