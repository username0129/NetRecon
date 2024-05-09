package request

import "github.com/gofrs/uuid/v5"

type UUIDRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
