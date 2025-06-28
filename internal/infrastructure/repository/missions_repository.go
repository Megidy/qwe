package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type MissionsRepository struct {
	pool *pgxpool.Pool
}

func NewMissionsRepository(pool *pgxpool.Pool) *MissionsRepository {
	return &MissionsRepository{
		pool: pool,
	}
}

func (r *MissionsRepository) CreateWithTargets(ctx context.Context, dto *model.MissionWithTargetsDTO) error {
	createMissionQuery := `
INSERT INTO missions (id, is_completed)
VALUES ($1, $2)	
`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin tx")
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error().Err(err).Msg("failed to rollback tx")
		}
	}()

	_, err = tx.Exec(
		ctx,
		createMissionQuery,
		dto.Mission.Id,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to execute create mission query: %w", err)
	}

	createTargetQuery := `
INSERT INTO targets(id, mission_id, name, country, notes, is_completed)
VALUES($1, $2, $3, $4, $5, $6)
	`

	for _, target := range dto.Targets {
		_, err = tx.Exec(
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
			return fmt.Errorf("failed to execute create target query: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (r *MissionsRepository) UpdateCatId(ctx context.Context, missionId, catId uuid.UUID) error {
	updateCatId := `
UPDATE missions SET cat_id=$1
WHERE id=$2
`
	_, err := r.pool.Exec(ctx, updateCatId, catId, missionId)
	if err != nil {
		return fmt.Errorf("failed to execute update cat id query: %w", err)
	}

	return nil
}

func (r *MissionsRepository) GetCatId(ctx context.Context, missionId uuid.UUID) (*uuid.UUID, error) {
	getCatIdQuery := `
SELECT cat_id 
FROM missions 
WHERE id=$1
	`

	var catId *uuid.UUID
	err := r.pool.QueryRow(ctx, getCatIdQuery, missionId).Scan(&catId)
	if err != nil {
		if isNoRows(err) {
			return nil, businesserrors.ErrMissionNotFound
		}
		return nil, fmt.Errorf("failed to query row: %w", err)
	}
	return catId, nil
}

func (r *MissionsRepository) GetCompletedStatusById(ctx context.Context, missionId uuid.UUID) (bool, error) {
	getStatusQuery := `
SELECT is_completed FROM missions
WHERE id=$1
	`
	var isCompleted bool
	err := r.pool.QueryRow(ctx, getStatusQuery, missionId).Scan(&isCompleted)
	if err != nil {
		return false, fmt.Errorf("failed to query row: %w", err)
	}

	return isCompleted, nil
}

func (r *MissionsRepository) Exists(ctx context.Context, missionId uuid.UUID) (bool, error) {
	existsQuery := `
SELECT EXISTS(
	SELECT 1 FROM missions
	WHERE id=$1
)
`
	var exists bool
	err := r.pool.QueryRow(ctx, existsQuery, missionId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to query if target exists: %w", err)
	}

	return exists, nil
}

func (r *MissionsRepository) UpdateMissionStatus(ctx context.Context, missionId uuid.UUID) error {
	updateTargetStatusQuery := `
UPDATE missions SET 
is_completed=true
WHERE id=$1
	`

	_, err := r.pool.Exec(ctx, updateTargetStatusQuery, missionId)
	if err != nil {
		return fmt.Errorf("failed to execute update mission status: %w", err)
	}

	return nil
}

func (r *MissionsRepository) GetMissionById(ctx context.Context, missionId uuid.UUID) (*model.MissionWithTargetsDTO, error) {
	query := `
SELECT
    m.id,
    m.cat_id,
    m.is_completed,
    m.created_at,
    m.updated_at,

    t.id,
    t.name,
    t.country,
    t.notes,
    t.is_completed,
    t.created_at,
    t.updated_at
FROM missions m
LEFT JOIN targets t ON t.mission_id = m.id
WHERE m.id = $1;
`

	rows, err := r.pool.Query(ctx, query, missionId)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var mission *model.MissionWithTargetsDTO
	targets := make([]*model.Target, 0)

	for rows.Next() {
		var (
			mID         uuid.UUID
			catID       *uuid.UUID
			isCompleted bool
			createdAt   time.Time
			updatedAt   time.Time

			targetID        *uuid.UUID
			targetName      *string
			targetCountry   *string
			targetNotes     *string
			targetCompleted *bool
			targetCreatedAt *time.Time
			targetUpdatedAt *time.Time
		)

		err = rows.Scan(
			&mID,
			&catID,
			&isCompleted,
			&createdAt,
			&updatedAt,
			&targetID,
			&targetName,
			&targetCountry,
			&targetNotes,
			&targetCompleted,
			&targetCreatedAt,
			&targetUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		if mission == nil {
			mission = &model.MissionWithTargetsDTO{
				Mission: &model.Mission{
					Id:          mID,
					CatId:       catID,
					IsCompleted: isCompleted,
					CreatedAt:   createdAt,
					UpdatedAt:   updatedAt,
				},
				Targets: []*model.Target{},
			}
		}
		if targetID != nil {
			target := &model.Target{
				Id:          *targetID,
				Name:        *targetName,
				Country:     *targetCountry,
				Notes:       *targetNotes,
				IsCompleted: *targetCompleted,
				CreatedAt:   *targetCreatedAt,
				UpdatedAt:   *targetUpdatedAt,
			}
			targets = append(targets, target)
		}
	}

	if mission == nil {
		return nil, businesserrors.ErrMissionNotFound
	}

	mission.Targets = targets

	return mission, nil
}

func (r *MissionsRepository) GetMissionsWithTargets(ctx context.Context, limit, offset int) ([]*model.MissionWithTargetsDTO, error) {
	getMissionsQuery := `
SELECT
    id,
    cat_id,
    is_completed,
    created_at,
    updated_at
FROM missions
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
`

	rows, err := r.pool.Query(ctx, getMissionsQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query missions: %w", err)
	}
	defer rows.Close()

	var missions []*model.MissionWithTargetsDTO
	missionIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var (
			mID         uuid.UUID
			catID       *uuid.UUID
			isCompleted bool
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(&mID, &catID, &isCompleted, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan mission: %w", err)
		}

		mission := &model.MissionWithTargetsDTO{
			Mission: &model.Mission{
				Id:          mID,
				CatId:       catID,
				IsCompleted: isCompleted,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			},
			Targets: []*model.Target{},
		}

		missions = append(missions, mission)
		missionIDs = append(missionIDs, mID)
	}

	if len(missionIDs) == 0 {
		return missions, nil
	}

	getTargetsQuery := `
SELECT
    id,
    mission_id,
    name,
    country,
    notes,
    is_completed,
    created_at,
    updated_at
FROM targets
WHERE mission_id = ANY($1);
`

	targetRows, err := r.pool.Query(ctx, getTargetsQuery, missionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to query targets: %w", err)
	}
	defer targetRows.Close()

	missionMap := make(map[uuid.UUID]*model.MissionWithTargetsDTO)
	for _, m := range missions {
		missionMap[m.Mission.Id] = m
	}

	for targetRows.Next() {
		var (
			tID         uuid.UUID
			missionID   uuid.UUID
			name        string
			country     string
			notes       string
			isCompleted bool
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := targetRows.Scan(&tID, &missionID, &name, &country, &notes, &isCompleted, &createdAt, &updatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan target: %w", err)
		}

		target := &model.Target{
			Id:          tID,
			Name:        name,
			Country:     country,
			Notes:       notes,
			IsCompleted: isCompleted,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		if mission, ok := missionMap[missionID]; ok {
			mission.Targets = append(mission.Targets, target)
		}
	}

	return missions, nil
}

func (r *MissionsRepository) Delete(ctx context.Context, missionId uuid.UUID) error {
	deleteQuery := `
DELETE FROM missions 
WHERE id=$1
	`

	_, err := r.pool.Exec(ctx, deleteQuery, missionId)
	if err != nil {
		return err
	}
	return nil
}
