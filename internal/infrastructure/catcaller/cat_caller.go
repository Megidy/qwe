package catcaller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Megidy/cats/internal/domain/model"
)

const (
	apiEndpoint = "https://api.thecatapi.com/v1/breeds"
)

type CatCaller struct {
}

func NewCatCaller() *CatCaller {
	return &CatCaller{}
}

func (c *CatCaller) GetCatBreeds(ctx context.Context) ([]*model.GetCatBreedResp, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiEndpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var breeds []*model.GetCatBreedResp
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	return breeds, nil
}
