package model

import (
	"time"

	"github.com/google/uuid"
)

type Target struct {
	Id        uuid.UUID
	MissionId uuid.UUID

	Name    string
	Country string
	Notes   string

	IsCompleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
