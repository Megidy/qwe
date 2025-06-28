package repository

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type CatsRepository interface {
	Create(ctx context.Context, cat *model.Cat) error
	UpdateSalary(ctx context.Context, catId uuid.UUID, salary float64) error
	Delete(ctx context.Context, catId uuid.UUID) error
	Exists(ctx context.Context, catId uuid.UUID) (bool, error)
	GetById(ctx context.Context, catId uuid.UUID) (*model.Cat, error)
	GetWithPagination(ctx context.Context, limit, offset int) ([]*model.Cat, error)
}
