package service

import (
	"context"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/Megidy/cats/internal/domain/repository"
	"github.com/google/uuid"
)

const (
	maxNumberOfTargets = 3
)

type TargetService struct {
	targetRepository  repository.TargetsRepository
	missionRepository repository.MissionsRepository
}

func NewTargetService(targetRepository repository.TargetsRepository, missionRepository repository.MissionsRepository) *TargetService {
	return &TargetService{
		targetRepository:  targetRepository,
		missionRepository: missionRepository,
	}
}

func (s *TargetService) Create(ctx context.Context, target *model.Target, missionId uuid.UUID) error {
	dto, err := s.missionRepository.GetMissionById(ctx, missionId)
	if err != nil {
		return err
	}

	if dto.Mission.IsCompleted {
		return businesserrors.ErrCannotAddTargetToCompletedMission
	}

	if len(dto.Targets) >= maxNumberOfTargets {
		return businesserrors.ErrMaxNumberOfTargetsExeeded
	}

	target.Id = uuid.New()
	target.MissionId = missionId

	return s.targetRepository.Create(ctx, target)
}

func (s *TargetService) UpdateStatus(ctx context.Context, targetId uuid.UUID) error {
	exists, err := s.targetRepository.Exists(ctx, targetId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrTargetNotExists
	}

	err = s.targetRepository.UpdateTargetStatus(ctx, targetId)
	if err != nil {
		return err
	}

	missionId, err := s.targetRepository.GetMissionIdByTargetId(ctx, targetId)
	if err != nil {
		return err
	}

	dto, err := s.missionRepository.GetMissionById(ctx, missionId)
	if err != nil {
		return err
	}

	targetsCount := len(dto.Targets)
	var targetsCompleted int
	for _, target := range dto.Targets {
		if target.IsCompleted {
			targetsCompleted++
		}
	}

	if targetsCount == targetsCompleted {
		return s.missionRepository.UpdateMissionStatus(ctx, missionId)
	}
	return nil
}

func (s *TargetService) UpdateNote(ctx context.Context, targetId uuid.UUID, note string) error {
	exists, err := s.targetRepository.Exists(ctx, targetId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrTargetNotExists
	}

	targetIsCompleted, err := s.targetRepository.GetCompletedStatusById(ctx, targetId)
	if err != nil {
		
		return err
	}
	if targetIsCompleted {
		return businesserrors.ErrTargetNotesCantBeUpdated
	}

	missionId, err := s.targetRepository.GetMissionIdByTargetId(ctx, targetId)
	if err != nil {
		return err
	}

	mission, err := s.missionRepository.GetMissionById(ctx, missionId)
	if err != nil {
		return err
	}

	if mission.Mission.IsCompleted {
		return businesserrors.ErrTargetNotesCantBeUpdated
	}

	return s.targetRepository.UpdateTargetNotes(ctx, targetId, note)
}

func (s *TargetService) Delete(ctx context.Context, targetId uuid.UUID) error {
	exists, err := s.targetRepository.Exists(ctx, targetId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrTargetNotExists
	}

	isCompleted, err := s.targetRepository.GetCompletedStatusById(ctx, targetId)
	if err != nil {
		return err
	}

	if isCompleted {
		return businesserrors.ErrTargetCanNotBeDeleted
	}
	return s.targetRepository.Delete(ctx, targetId)
}
