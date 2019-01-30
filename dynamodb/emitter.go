package dynamodb

import (
	"context"
	"fmt"
	"wumber"
)

type emitClient struct {
	emitter wumber.UserEmitter
	next    wumber.Storage
}

// WrapWithEmitter wraps an Storage and start emits.
func WrapWithEmitter(emitter wumber.UserEmitter, c wumber.Storage) wumber.Storage {
	return &emitClient{emitter, c}
}

func (c *emitClient) Register(ctx context.Context, u wumber.User) (user wumber.User, err error) {
	defer func() {
		if err != nil {
			return
		}

		err2 := c.emitter.EmitCreate(ctx, user)
		fmt.Println(err2)
	}()

	return c.next.Register(ctx, u)
}

func (c *emitClient) CreateWorkspace(ctx context.Context, w wumber.Workspace) (wumber.WorkspaceID, error) {
	return c.next.CreateWorkspace(ctx, w)
}
