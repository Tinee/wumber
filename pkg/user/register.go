package user

import (
	"context"
	"strings"
	"time"
	"wumber"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

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
	jwt, err := s.jwtService.Extract(user)
	if err != nil {
		return "", errors.Wrap(err, "error extracting the token from the user")
	}

	return jwt, nil
}

func (s *loggingService) Register(ctx context.Context, input RegisterUserInput) (jwt wumber.JWT, err error) {
	defer s.logger.Flush()
	s.logger.Debug(ctx, "Calling Register",
		"input", input,
	)
	defer func(begin time.Time) {
		if err != nil {
			s.logger.Error(ctx, "Failed to call Register.",
				"input", input,
				"err", err,
				"took", time.Since(begin),
			)
			return
		}

		s.logger.Debug(ctx, "Called Register",
			"input", input,
			"output", jwt,
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.next.Register(ctx, input)
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
		FirstName: strings.Title(strings.ToLower(strings.TrimSpace(u.FirstName))),
		LastName:  strings.Title(strings.ToLower(strings.TrimSpace(u.LastName))),
		Email:     strings.ToLower(strings.TrimSpace(u.Email)),
		Password:  u.Password,
	}, nil
}
