package user

import (
	"wumber"

	jwt "github.com/dgrijalva/jwt-go"
)

// JWT is the json-web-token we're going to give back when the user either
// 1: Login
// 2: Register
type JWT string

func (s *service) extractJWT(u wumber.User) (JWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: 15000,
		Id:        string(u.ID),
	})
	ss, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return JWT(ss), nil
}
