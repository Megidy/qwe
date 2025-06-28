package http

import (
	"errors"
	"net/http"

	businesserrors "github.com/Megidy/cats/internal/domain/errors"
	"github.com/Megidy/cats/internal/domain/model"
	"github.com/Megidy/cats/internal/domain/service"
	"github.com/Megidy/cats/internal/interface/http/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

const (
	missionIDParam = "mission_id"
)

type MissionsHandler struct {
	missionService service.MissionService
}

func NewMissionsHandler(missionService service.MissionService) *MissionsHandler {
	return &MissionsHandler{
		missionService: missionService,
	}
}

func (h *MissionsHandler) Create(ctx echo.Context) error {
	var req dto.CreateMissionDTO

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	dtoObj := model.MissionWithTargetsDTO{
		Mission: &model.Mission{},
	}

	var targets = make([]*model.Target, 0, len(req.Targets))

	for _, target := range req.Targets {
		targets = append(targets, &model.Target{
			Name:    target.Name,
			Country: target.Country,
			Notes:   target.Notes,
		})
	}
	dtoObj.Targets = targets

	id, err := h.missionService.Create(ctx.Request().Context(), &dtoObj)
	if err != nil {
		log.Error().Err(err).Msg("failed to create mission")
		switch {
		case errors.Is(err, businesserrors.ErrMaxNumberOfTargetsExeeded):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	resp := dto.CreateResponse{
		Id: id.String(),
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully created mission", http.StatusOK, resp))
}

func (h *MissionsHandler) Delete(ctx echo.Context) error {
	missionId := getFromParam(ctx, missionIDParam)

	uuid, err := uuid.Parse(missionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.missionService.Delete(ctx.Request().Context(), uuid)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete mission")
		switch {
		case errors.Is(err, businesserrors.ErrCatIsAssigned):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully deleted mission", http.StatusOK, nil))
}

func (h *MissionsHandler) UpdateStatus(ctx echo.Context) error {
	missionId := getFromParam(ctx, missionIDParam)

	uuid, err := uuid.Parse(missionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.missionService.UpdateStatus(ctx.Request().Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrMissionNotFound):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		log.Error().Err(err).Msg("failed to update mission")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully updated mission", http.StatusOK, nil))

}
func (h *MissionsHandler) AssignCat(ctx echo.Context) error {
	missionId := getFromParam(ctx, missionIDParam)

	var req dto.AssigneCatDTO
	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	missionUUID, err := uuid.Parse(missionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	catId, err := uuid.Parse(req.CatId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.missionService.AssignCat(ctx.Request().Context(), missionUUID, catId)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrMissionNotFound) ||
			errors.Is(err, businesserrors.ErrCatNotFound):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}

		log.Error().Err(err).Msg("failed to update mission")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully assigned cat", http.StatusOK, nil))
}
func (h *MissionsHandler) GetMissions(ctx echo.Context) error {
	limit, offset, err := getAndValidateLimitAndPage(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	missions, err := h.missionService.GetMissions(ctx.Request().Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to get missions")
		return ctx.JSON(http.StatusBadRequest, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully retrieved missions", http.StatusOK, missions))
}
func (h *MissionsHandler) GetMissionById(ctx echo.Context) error {
	missionId := getFromParam(ctx, missionIDParam)

	missionUUID, err := uuid.Parse(missionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	mission, err := h.missionService.GetMissionById(ctx.Request().Context(), missionUUID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get mission")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully retrieved mission", http.StatusOK, mission))
}
