package dto

import "net/http"

type CreateCatDTO struct {
	Name   string  `json:"name"`
	Breed  string  `json:"breed"`
	Salary float64 `json:"salary"`
}

type CreateResponse struct {
	Id string `json:"id"`
}

type GeneralResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
	Data       any    `json:"data"`
}

type UpdateSalaryReq struct {
	Salary float64 `json:"salary"`
}

type CreateMissionDTO struct {
	Targets []TargetInputDTO `json:"targets"`
}

type TargetInputDTO struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Notes   string `json:"notes,omitempty"`
}

type AssigneCatDTO struct {
	CatId string `json:"cat_id"`
}

type UpdateNotesDTO struct {
	Notes string `json:"notes"`
}

func NewBadRequest(err error) GeneralResponse {
	return GeneralResponse{
		Success:    false,
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
		Data:       nil,
	}
}

func NewInternalServerError(err error) GeneralResponse {
	return GeneralResponse{
		Success:    false,
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
		Data:       nil,
	}
}

func NewSuccessGeneralResponse(message string, statusCode int, data any) GeneralResponse {
	return GeneralResponse{
		Success:    true,
		Message:    message,
		StatusCode: statusCode,
		Data:       data,
	}
}
