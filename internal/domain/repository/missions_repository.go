package repository

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type MissionsRepository interface {
	CreateWithTargets(ctx context.Context, dto *model.MissionWithTargetsDTO) error
	UpdateCatId(ctx context.Context, missionId, catId uuid.UUID) error
	GetCatId(ctx context.Context, missionId uuid.UUID) (*uuid.UUID, error)
	GetCompletedStatusById(ctx context.Context, missionId uuid.UUID) (bool, error)
	Exists(ctx context.Context, missionId uuid.UUID) (bool, error)
	UpdateMissionStatus(ctx context.Context, missionId uuid.UUID) error
	GetMissionById(ctx context.Context, missionId uuid.UUID) (*model.MissionWithTargetsDTO, error)
	GetMissionsWithTargets(ctx context.Context, limit, offset int) ([]*model.MissionWithTargetsDTO, error)
	Delete(ctx context.Context, missionId uuid.UUID) error
}
