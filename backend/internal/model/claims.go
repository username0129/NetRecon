package model

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UUID        uuid.UUID
	ID          uint
	Username    string
	AuthorityId uint
	BufferTime  int64
	jwt.RegisteredClaims
}
