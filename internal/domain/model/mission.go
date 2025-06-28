package model

import (
	"time"

	"github.com/google/uuid"
)

type Mission struct {
	Id    uuid.UUID
	CatId *uuid.UUID

	IsCompleted bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
