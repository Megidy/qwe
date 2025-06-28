package repository

import (
	"context"
	"fmt"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatsRepository struct {
	pool *pgxpool.Pool
}

func NewCatsRepository(pool *pgxpool.Pool) *CatsRepository {
	return &CatsRepository{
		pool: pool,
	}
}

func (r *CatsRepository) Create(ctx context.Context, cat *model.Cat) error {
	createCatQuery := `
	INSERT INTO cats (id, name, started_at, on_mission, breed, salary)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(
		ctx,
		createCatQuery,
		cat.Id,
		cat.Name,
		cat.StartedAt,
		cat.OnMission,
		cat.Breed,
		cat.Salary,
	)

	if err != nil {
		return fmt.Errorf("failed to execute create cat query: %w", err)
	}

	return nil
}

func (r *CatsRepository) UpdateSalary(ctx context.Context, catId uuid.UUID, salary float64) error {
	updateSalaryQuery := `
		UPDATE cats SET salary = $1, updated_at = NOW() WHERE id = $2
	`

	_, err := r.pool.Exec(ctx, updateSalaryQuery,
		salary,
		catId,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update salary query: %w", err)
	}

	return nil
}

func (r *CatsRepository) Delete(ctx context.Context, catId uuid.UUID) error {
	deleteCatQuery := `
		DELETE FROM cats WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, deleteCatQuery, catId)
	if err != nil {
		return fmt.Errorf("failed to execute delete cat query: %w", err)
	}

	return nil
}

func (r *CatsRepository) Exists(ctx context.Context, catId uuid.UUID) (bool, error) {
	existsQuery := `
SELECT EXISTS(
	SELECT 1 FROM cats
	WHERE id=$1
)
`
	var exists bool
	err := r.pool.QueryRow(ctx, existsQuery, catId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to query if cat exists: %w", err)
	}

	return exists, nil
}

func (r *CatsRepository) GetById(ctx context.Context, catId uuid.UUID) (*model.Cat, error) {
	getCatByIdQuery := `
		SELECT id, name, started_at, on_mission, breed, salary, created_at, updated_at
		FROM cats WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, getCatByIdQuery, catId)

	var cat model.Cat
	err := row.Scan(
		&cat.Id,
		&cat.Name,
		&cat.StartedAt,
		&cat.OnMission,
		&cat.Breed,
		&cat.Salary,
		&cat.CreatedAt,
		&cat.UpdatedAt,
	)
	if err != nil {
		if isNoRows(err) {
			return nil, businesserrors.ErrCatNotFound
		}
		return nil, fmt.Errorf("failed to scan cat by id: %w", err)
	}

	return &cat, nil
}

func (r *CatsRepository) GetWithPagination(ctx context.Context, limit, offset int) ([]*model.Cat, error) {
	getCatsQuery := `
		SELECT id, name, started_at, on_mission, breed, salary, created_at, updated_at
		FROM cats
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, getCatsQuery, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get cats query: %w", err)
	}
	defer rows.Close()

	var cats []*model.Cat
	for rows.Next() {
		var cat model.Cat
		err := rows.Scan(
			&cat.Id,
			&cat.Name,
			&cat.StartedAt,
			&cat.OnMission,
			&cat.Breed,
			&cat.Salary,
			&cat.CreatedAt,
			&cat.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cat row: %w", err)
		}
		cats = append(cats, &cat)
	}

	return cats, nil
}
