package wumber

import (
	"context"
	"time"
)

// User is the domain object of users.
type User struct {
	ID        UserID    `dynamodbav:"Id"`
	Email     string    `dynamodbav:"Email"`
	Password  string    `dynamodbav:"Password"`
	FirstName string    `dynamodbav:"FirstName"`
	LastName  string    `dynamodbav:"LastName"`
	Created   time.Time `dynamodbav:"Created"`
}

// UserID is an ID that, with it we can identify users.
type UserID string

// UserRepository is an interface that can store users.
type UserRepository interface {
	Register(ctx context.Context, user User) (User, error)
}

// UserEmitter is a domain interface that describes something that will emit events for user purposes.
type UserEmitter interface {
	EmitCreate(ctx context.Context, new User) error
}
