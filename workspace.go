package wumber

import (
	"context"
	"time"
)

// Workspace represents a workspace where user in this workspace can keep track of their scores.
type Workspace struct {
	ID      WorkspaceID `dynamodbav:"ID"`
	Name    string      `dynamodbav:"Name"`
	Created time.Time   `dynamodbav:"Created"`
}

type WorkspaceID string

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, name string) (WorkspaceID, error)
}
