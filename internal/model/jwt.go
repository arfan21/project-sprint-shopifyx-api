package model

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	// UserID for middleware purpose
	UserID uuid.UUID `json:"-"`

	Name string `json:"name"`
	jwt.RegisteredClaims
}
