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
	targetParam = "target_id"
)

type TargetsHandler struct {
	targetSertvice service.TargetService
}

func NewTargetsHandler(targetSertvice service.TargetService) *TargetsHandler {
	return &TargetsHandler{
		targetSertvice: targetSertvice,
	}
}

func (h *TargetsHandler) Create(ctx echo.Context) error {
	missionId := getFromParam(ctx, missionIDParam)

	uuid, err := uuid.Parse(missionId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	var req dto.TargetInputDTO

	err = ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	target := &model.Target{
		Name:    req.Name,
		Country: req.Country,
		Notes:   req.Notes,
	}
	err = h.targetSertvice.Create(ctx.Request().Context(), target, uuid)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrCannotAddTargetToCompletedMission) ||
			errors.Is(err, businesserrors.ErrMaxNumberOfTargetsExeeded):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		log.Error().Err(err).Msg("failed to create target")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully created target", http.StatusOK, nil))

}

func (h *TargetsHandler) UpdateStatus(ctx echo.Context) error {
	targetId := getFromParam(ctx, targetParam)

	uuid, err := uuid.Parse(targetId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}
	err = h.targetSertvice.UpdateStatus(ctx.Request().Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrTargetNotExists):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		log.Error().Err(err).Msg("failed to update target")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully updated target", http.StatusOK, nil))
}

func (h *TargetsHandler) UpdateNote(ctx echo.Context) error {
	targetId := getFromParam(ctx, targetParam)

	uuid, err := uuid.Parse(targetId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	var req dto.UpdateNotesDTO

	err = ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.targetSertvice.UpdateNote(ctx.Request().Context(), uuid, req.Notes)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrTargetNotExists) ||
			errors.Is(err, businesserrors.ErrTargetNotesCantBeUpdated):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		log.Error().Err(err).Msg("failed to update target notes")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully update target notes", http.StatusOK, nil))
}

func (h *TargetsHandler) Delete(ctx echo.Context) error {
	targetId := getFromParam(ctx, targetParam)

	uuid, err := uuid.Parse(targetId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.targetSertvice.Delete(ctx.Request().Context(), uuid)
	if err != nil {
		switch {
		case errors.Is(err, businesserrors.ErrTargetNotExists) ||
			errors.Is(err, businesserrors.ErrTargetCanNotBeDeleted):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		log.Error().Err(err).Msg("failed to delete target")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully delete target", http.StatusOK, nil))

}
