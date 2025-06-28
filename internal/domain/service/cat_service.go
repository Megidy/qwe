package service

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
	"github.com/google/uuid"
)

type CatService interface {
	Create(ctx context.Context, cat *model.Cat) (uuid.UUID,error)
	UpdateSalary(ctx context.Context, catId uuid.UUID, salary float64) error
	Delete(ctx context.Context, catId uuid.UUID) error
	GetById(ctx context.Context, catId uuid.UUID) (*model.Cat, error)
	GetWithPagination(ctx context.Context, limit, offset int) ([]*model.Cat, error)
}
