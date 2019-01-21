package wumber

import (
	"errors"
)

// Common errors in the domain.
var (
	ErrRegisterUserEmailExists     = errors.New("email exists")
	ErrCreatingWorkspaceNameExists = errors.New("workspace exists")
)
