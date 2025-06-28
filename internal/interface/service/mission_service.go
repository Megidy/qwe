package service

import (
	"context"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/Megidy/cats/internal/domain/repository"
	"github.com/google/uuid"
)

type MissionService struct {
	missionRepository repository.MissionsRepository
	catsRepository    repository.CatsRepository
}

func NewMissionService(missionRepository repository.MissionsRepository, catsRepository repository.CatsRepository) *MissionService {
	return &MissionService{
		missionRepository: missionRepository,
		catsRepository:    catsRepository,
	}
}

func (s *MissionService) Create(ctx context.Context, dto *model.MissionWithTargetsDTO) (uuid.UUID, error) {
	dto.Mission.Id = uuid.New()
	if len(dto.Targets) >= maxNumberOfTargets {
		return uuid.UUID{}, businesserrors.ErrMaxNumberOfTargetsExeeded
	}

	for _, target := range dto.Targets {
		target.Id = uuid.New()
		target.MissionId = dto.Mission.Id
	}

	return dto.Mission.Id, s.missionRepository.CreateWithTargets(ctx, dto)
}

func (s *MissionService) Delete(ctx context.Context, missionId uuid.UUID) error {
	catId, err := s.missionRepository.GetCatId(ctx, missionId)
	if err != nil {
		return err
	}

	if catId != nil {
		return businesserrors.ErrCatIsAssigned
	}

	return s.missionRepository.Delete(ctx, missionId)
}

func (s *MissionService) UpdateStatus(ctx context.Context, missionId uuid.UUID) error {
	exists, err := s.missionRepository.Exists(ctx, missionId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrMissionNotFound
	}

	return s.missionRepository.UpdateMissionStatus(ctx, missionId)
}

func (s *MissionService) AssignCat(ctx context.Context, missionId uuid.UUID, catId uuid.UUID) error {
	exists, err := s.missionRepository.Exists(ctx, missionId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrMissionNotFound
	}

	exists, err = s.catsRepository.Exists(ctx, catId)
	if err != nil {
		return err
	}
	
	if !exists {
		return businesserrors.ErrCatNotFound
	}

	return s.missionRepository.UpdateCatId(ctx, missionId, catId)
}

func (s *MissionService) GetMissions(ctx context.Context, limit, offset int) ([]*model.MissionWithTargetsDTO, error) {
	return s.missionRepository.GetMissionsWithTargets(ctx, limit, offset)
}

func (s *MissionService) GetMissionById(ctx context.Context, missionId uuid.UUID) (*model.MissionWithTargetsDTO, error) {
	return s.missionRepository.GetMissionById(ctx, missionId)
}
