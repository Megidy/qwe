package service

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type MissionService interface {
	Create(ctx context.Context, dto *model.MissionWithTargetsDTO) (uuid.UUID,error)
	Delete(ctx context.Context, missionId uuid.UUID) error
	UpdateStatus(ctx context.Context, missionId uuid.UUID) error
	AssignCat(ctx context.Context, missionId uuid.UUID, catId uuid.UUID) error
	GetMissions(ctx context.Context, limit, offset int) ([]*model.MissionWithTargetsDTO, error)
	GetMissionById(ctx context.Context, missionId uuid.UUID) (*model.MissionWithTargetsDTO, error)
}
