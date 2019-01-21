package wumber

import (
	"context"
	"time"
)

// Workspace represents a workspace where user in this workspace can keep track of their scores.
type Workspace struct {
	ID      WorkspaceID `dynamodbav:"Id"`
	OwnerID string      `dynamodbav:"OwnerId"`
	Name    string      `dynamodbav:"Name"`
	Created time.Time   `dynamodbav:"Created"`
}

// WorkspaceID is an ID that identifies workspaces.
type WorkspaceID string

// WorkspaceRepository is an interface that can store users.
type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, w Workspace) (WorkspaceID, error)
}
