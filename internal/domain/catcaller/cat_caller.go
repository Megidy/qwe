package catcaller

import (
	"context"

	"github.com/Megidy/cats/internal/domain/model"
)

type CatCaller interface {
	GetCatBreeds(ctx context.Context) ([]*model.GetCatBreedResp,error)
}
