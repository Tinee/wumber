package jwt

import (
	"wumber"

	jwt "github.com/dgrijalva/jwt-go"
)

// Service represents the container for the JWTService.
type Service struct {
	secret wumber.JWTSecret
}

// New instantiate a JWTService.
func New(s wumber.JWTSecret) *Service {
	return &Service{s}
}

// Extract takes a User and makes it a JWT.
func (f *Service) Extract(u wumber.User) (wumber.JWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: 15000,
		Id:        string(u.ID),
	})
	ss, err := token.SignedString([]byte(f.secret))
	if err != nil {
		return "", err
	}

	return wumber.JWT(ss), nil
}
