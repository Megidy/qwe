package repository

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type TargetsRepository interface {
	Create(ctx context.Context, target *model.Target) error
	GetCompletedStatusById(ctx context.Context, targetId uuid.UUID) (bool, error)
	Exists(ctx context.Context, targetId uuid.UUID) (bool, error)
	GetAmountOfTargetsOfMission(ctx context.Context, missionId uuid.UUID) (int, error)
	UpdateTargetStatus(ctx context.Context, targetId uuid.UUID) error
	UpdateTargetNotes(ctx context.Context, targetId uuid.UUID, notes string) error
	Delete(ctx context.Context, targetId uuid.UUID) error
	GetMissionIdByTargetId(ctx context.Context, targetId uuid.UUID) (uuid.UUID, error)
}
