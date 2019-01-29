package wumber

// JWT represents the JWT in the domain.
type JWT string

// JWTSecret is a secret we're using to Extract and Validate users.
type JWTSecret string

// JWTService is a representation how a JWTService.
type JWTService interface {
	Extract(User) (JWT, error)
}
