package user

import (
	"context"
	"strings"
	"wumber"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

// RegisterUserInput is the struct we need to create a user.
type RegisterUserInput struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
}

func (s *service) Register(ctx context.Context, input RegisterUserInput) (wumber.JWT, error) {
	u, err := input.toUser()
	if err != nil {
		return "", errors.Wrap(err, "error converting the input to an user.")
	}
	bs, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	u.Password = string(bs)

	user, err := s.userRepo.Register(ctx, u)
	if err != nil {
		return "", errors.Wrap(err, "error when register the user")
	}
	jwt, err := s.jwtService.Extract(user, s.secret)
	if err != nil {
		return "", errors.Wrap(err, "error extracting the token from the user")
	}

	return jwt, nil
}

// toUser validates the user, normalizes some fields and converts it to a domain User.
func (u RegisterUserInput) toUser() (wumber.User, error) {
	err := validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 50)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 50)),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 128)),
	)
	if err != nil {
		return wumber.User{}, err
	}

	return wumber.User{
		FirstName: strings.ToTitle(strings.TrimSpace(u.FirstName)),
		LastName:  strings.ToTitle(strings.TrimSpace(u.LastName)),
		Email:     strings.ToLower(u.Email),
		Password:  u.Password,
	}, nil
}
