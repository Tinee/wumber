package jwt

import (
	"wumber"

	jwt "github.com/dgrijalva/jwt-go"
)

// Extract takes a User and makes it a JWT.
func Extract(u wumber.User, s wumber.JWTSecret) (wumber.JWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: 15000,
		Id:        string(u.ID),
	})
	ss, err := token.SignedString([]byte(s))
	if err != nil {
		return "", err
	}

	return wumber.JWT(ss), nil
}
