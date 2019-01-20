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

type WorkspaceID string

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, name, accountID string) (WorkspaceID, error)
}
