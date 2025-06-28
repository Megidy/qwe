package repository

import (
	"context"
	"fmt"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TargetsRepository struct {
	pool *pgxpool.Pool
}

func NewTargetsRepository(pool *pgxpool.Pool) *TargetsRepository {
	return &TargetsRepository{
		pool: pool,
	}
}

func (r *TargetsRepository) Create(ctx context.Context, target *model.Target) error {
	createTargetQuery := `
INSERT INTO targets(id, mission_id, name, country, notes, is_completed)
VALUES($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		createTargetQuery,
		target.Id,
		target.MissionId,
		target.Name,
		target.Country,
		target.Notes,
		target.IsCompleted,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *TargetsRepository) GetCompletedStatusById(ctx context.Context, targetId uuid.UUID) (bool, error) {
	getStatusQuery := `
SELECT is_completed FROM targets
WHERE id=$1
	`
	var isCompleted bool
	err := r.pool.QueryRow(ctx, getStatusQuery, targetId).Scan(&isCompleted)
	if err != nil {
		if isNoRows(err) {
			return false, businesserrors.ErrTargetNotExists
		}
		return false, fmt.Errorf("failed to query row: %w", err)
	}

	return isCompleted, nil
}

func (r *TargetsRepository) Delete(ctx context.Context, targetId uuid.UUID) error {
	deleteTargetQuery := `
DELETE FROM targets
WHERE id=$1
	`

	_, err := r.pool.Exec(ctx, deleteTargetQuery, targetId)
	if err != nil {
		return fmt.Errorf("failed to execute delete target query: %w", err)
	}

	return nil
}

func (r *TargetsRepository) Exists(ctx context.Context, targetId uuid.UUID) (bool, error) {
	existsQuery := `
SELECT EXISTS(
	SELECT 1 FROM targets
	WHERE id=$1
)
	`
	var exists bool
	err := r.pool.QueryRow(ctx, existsQuery, targetId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to query if target exists: %w", err)
	}

	return exists, nil
}

func (r *TargetsRepository) GetAmountOfTargetsOfMission(ctx context.Context, missionId uuid.UUID) (int, error) {
	getAmountQuery := `
SELECT COUNT(*) FROM targets
WHERE mission_id=$1
	`

	var count int

	err := r.pool.QueryRow(ctx, getAmountQuery, missionId).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to query amount of targets: %w", err)
	}

	return count, nil
}

func (r *TargetsRepository) UpdateTargetStatus(ctx context.Context, targetId uuid.UUID) error {
	updateTargetStatusQuery := `
UPDATE targets SET 
is_completed=true
WHERE id=$1
	`

	_, err := r.pool.Exec(ctx, updateTargetStatusQuery, targetId)
	if err != nil {
		return fmt.Errorf("failed to execute update target status query: %w", err)
	}

	return nil
}

func (r *TargetsRepository) UpdateTargetNotes(ctx context.Context, targetId uuid.UUID, notes string) error {
	updateQuery := `
UPDATE targets SET notes=$1
WHERE id=$2
	`

	_, err := r.pool.Exec(ctx, updateQuery, notes, targetId)
	if err != nil {
		return fmt.Errorf("failed to execute update target notes query: %w", err)
	}

	return nil

}

func (r *TargetsRepository) GetMissionIdByTargetId(ctx context.Context, targetId uuid.UUID) (uuid.UUID, error) {
	getMissionIdQuery := `
SELECT mission_id
FROM targets
WHERE id=$1
	`

	var id uuid.UUID

	err := r.pool.QueryRow(ctx, getMissionIdQuery, targetId).Scan(&id)
	if err != nil {
		if isNoRows(err) {
			return uuid.UUID{}, businesserrors.ErrTargetNotExists
		}
		return uuid.UUID{}, fmt.Errorf("failed to query amount of targets: %w", err)
	}

	return id, nil

}
