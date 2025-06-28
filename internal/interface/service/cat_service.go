package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Megidy/cats/internal/domain/catcaller"
	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/Megidy/cats/internal/domain/repository"
	"github.com/google/uuid"
)

type CatService struct {
	catRepository repository.CatsRepository
	catCaller     catcaller.CatCaller
}

func NewCatService(catRepository repository.CatsRepository, catCaller catcaller.CatCaller) *CatService {
	return &CatService{
		catRepository: catRepository,
		catCaller:     catCaller,
	}
}

func (s *CatService) validateBreeds(breeds []*model.GetCatBreedResp, breed string) error {
	for _, b := range breeds {
		if b.Breed == breed {
			return nil
		}
	}
	return businesserrors.ErrBreedIsNotValid
}

func (s *CatService) Create(ctx context.Context, cat *model.Cat) (uuid.UUID, error) {
	breeds, err := s.catCaller.GetCatBreeds(ctx)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to get breeds: %w", err)
	}

	err = s.validateBreeds(breeds, cat.Breed)
	if err != nil {
		return uuid.UUID{}, err
	}

	cat.Id = uuid.New()
	cat.StartedAt = time.Now()

	return cat.Id, s.catRepository.Create(ctx, cat)
}

func (s *CatService) UpdateSalary(ctx context.Context, catId uuid.UUID, salary float64) error {
	exists, err := s.catRepository.Exists(ctx, catId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrCatNotFound
	}

	return s.catRepository.UpdateSalary(ctx, catId, salary)
}

func (s *CatService) Delete(ctx context.Context, catId uuid.UUID) error {
	exists, err := s.catRepository.Exists(ctx, catId)
	if err != nil {
		return err
	}

	if !exists {
		return businesserrors.ErrCatNotFound
	}

	return s.catRepository.Delete(ctx, catId)
}

func (s *CatService) GetById(ctx context.Context, catId uuid.UUID) (*model.Cat, error) {
	return s.catRepository.GetById(ctx, catId)
}

func (s *CatService) GetWithPagination(ctx context.Context, limit, offset int) ([]*model.Cat, error) {
	return s.catRepository.GetWithPagination(ctx, limit, offset)
}
