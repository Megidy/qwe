package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func checkHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]string{
		"health": "alive",
	})
}
