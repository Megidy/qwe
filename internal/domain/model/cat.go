package model

import (
	"time"

	"github.com/google/uuid"
)

type Cat struct {
	Id        uuid.UUID
	Name      string
	StartedAt time.Time

	OnMission bool

	Breed  string
	Salary float64

	CreatedAt time.Time
	UpdatedAt time.Time
}
