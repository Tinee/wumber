package user

import (
	"context"
	"wumber"
)

type Service interface {
	Register(context.Context, RegisterUserInput) (wumber.JWT, error)
}

type service struct {
	userRepo   wumber.UserRepository
	jwtService wumber.JWTService
}

func NewService(r wumber.UserRepository, jwt wumber.JWTService) Service {
	return &service{
		userRepo:   r,
		jwtService: jwt,
	}
}
