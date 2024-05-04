package model

import (
	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UUID        uuid.UUID
	Username    string
	AuthorityId string
	BufferTime  int64
	jwt.RegisteredClaims
}
