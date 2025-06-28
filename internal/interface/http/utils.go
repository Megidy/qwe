package http

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	limitQueryParam = "limit"
	pageQueryParam  = "page"
)

func getAndValidateLimitAndPage(ctx echo.Context) (int, int, error) {

	limitParam := getQueryParam(ctx, limitQueryParam)
	pageParam := getQueryParam(ctx, pageQueryParam)

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse limit param: %w", err)
	}

	if limitParam == "" || pageParam == "" {
		return 0, 0, fmt.Errorf("invalid request query params: requires limit, page params")
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse page param: %w", err)
	}

	if page <= 0 || limit <= 0 {
		return 0, 0, fmt.Errorf("page and limit must be positive numbers: page: %d, limit: %d", page, limit)
	}

	offset := (page - 1) * limit
	return limit, offset, nil

}

func getQueryParam(ctx echo.Context, param string) string {
	return ctx.QueryParam(param)
}

func getFromParam(ctx echo.Context, param string) string {
	return ctx.Param(param)
}
