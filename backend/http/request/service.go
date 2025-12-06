package request

import "github.com/google/uuid"

type ServiceCreate struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
