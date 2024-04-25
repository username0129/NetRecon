package request

import "github.com/gofrs/uuid/v5"

type CancelTaskRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
