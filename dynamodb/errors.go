package dynamodb

import (
	"errors"
)

// Common errors that can happen when we interact with DynamoDB
var (
	ErrWorkspaceNameExists = errors.New("workspace is taken")
	ErrUserEmailExists     = errors.New("email is taken")
	ErrUnexpectedCause     = errors.New("unexpected error")
)
