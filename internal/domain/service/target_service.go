package service

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type TargetService interface {
	Create(ctx context.Context, target *model.Target, missionId uuid.UUID) error
	UpdateStatus(ctx context.Context, targetId uuid.UUID) error
	UpdateNote(ctx context.Context, targetId uuid.UUID, note string) error
	Delete(ctx context.Context, targetId uuid.UUID) error
}
