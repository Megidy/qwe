package model

type MissionWithTargetsDTO struct {
	Mission *Mission
	Targets []*Target
}

type GetCatBreedResp struct {
	Breed string `json:"name"`
}
