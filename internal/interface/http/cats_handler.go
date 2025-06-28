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
	catIdParam = "cat_id"
)

type CatsHandler struct {
	catService service.CatService
}

func NewCatsHandler(catService service.CatService) *CatsHandler {
	return &CatsHandler{
		catService: catService,
	}
}

func (h *CatsHandler) Create(ctx echo.Context) error {
	var req dto.CreateCatDTO

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	cat := model.Cat{
		Name:   req.Name,
		Breed:  req.Breed,
		Salary: req.Salary,
	}

	id, err := h.catService.Create(ctx.Request().Context(), &cat)
	if err != nil {
		log.Error().Err(err).Msg("failed to create cat")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}
	resp := dto.CreateResponse{
		Id: id.String(),
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully created cat", http.StatusOK, resp))
}

func (h *CatsHandler) UpdateSalary(ctx echo.Context) error {
	catId := getFromParam(ctx, catIdParam)

	var req dto.UpdateSalaryReq

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	uuid, err := uuid.Parse(catId)
	if err != nil {

		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.catService.UpdateSalary(ctx.Request().Context(), uuid, req.Salary)
	if err != nil {
		log.Error().Err(err).Msg("failed to update cat salary")
		switch {
		case errors.Is(err, businesserrors.ErrCatNotFound):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully updated salry", http.StatusOK, nil))
}

func (h *CatsHandler) Delete(ctx echo.Context) error {
	catId := getFromParam(ctx, catIdParam)

	uuid, err := uuid.Parse(catId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	err = h.catService.Delete(ctx.Request().Context(), uuid)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete cat")
		switch {
		case errors.Is(err, businesserrors.ErrCatNotFound):
			return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
		}
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully deleted cat", http.StatusOK, nil))
}

func (h *CatsHandler) GetById(ctx echo.Context) error {
	catId := getFromParam(ctx, catIdParam)

	uuid, err := uuid.Parse(catId)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	cat, err := h.catService.GetById(ctx.Request().Context(), uuid)
	if err != nil {
		log.Error().Err(err).Msg("failed to get cat by id")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully retrieved cat", http.StatusOK, cat))
}

func (h *CatsHandler) GetAll(ctx echo.Context) error {

	limit, offset, err := getAndValidateLimitAndPage(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, dto.NewBadRequest(err))
	}

	cats, err := h.catService.GetWithPagination(ctx.Request().Context(), limit, offset)
	if err != nil {
		log.Error().Err(err).Msg("failed to get cats")
		return ctx.JSON(http.StatusInternalServerError, dto.NewInternalServerError(err))
	}

	return ctx.JSON(http.StatusOK, dto.NewSuccessGeneralResponse("successfully retrieved cat", http.StatusOK, cats))
}
