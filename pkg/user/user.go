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
	secret     wumber.JWTSecret
}

func NewService(r wumber.UserRepository, jwt wumber.JWTService, secret wumber.JWTSecret) Service {
	return &service{
		userRepo:   r,
		secret:     secret,
		jwtService: jwt,
	}
}
