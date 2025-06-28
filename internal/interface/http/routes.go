package http

import (
	"github.com/Megidy/cats/internal/interface/middleware"
	"github.com/labstack/echo/v4"
)

const (
	apiVersion    = "/v1"
	catsGroup     = "/cats"
	missionsGroup = "/missions"
	targetsGroup  = "/targets"
)

type Router struct {
	echo            *echo.Echo
	catsHandler     *CatsHandler
	missionsHandler *MissionsHandler
	targetsHandler  *TargetsHandler
}

func NewRouter(
	echo *echo.Echo,
	catsHandler *CatsHandler,
	missionsHandler *MissionsHandler,
	targetsHandler *TargetsHandler,
) *Router {
	return &Router{
		echo:            echo,
		catsHandler:     catsHandler,
		missionsHandler: missionsHandler,
		targetsHandler:  targetsHandler,
	}
}

func (r *Router) RegisterRoutes() {
	api := r.echo.Group(apiVersion, middleware.WithRequestResponseLogger())

	api.GET("/health", checkHealth)
	cats := api.Group(catsGroup)
	missions := api.Group(missionsGroup)
	targets := api.Group(targetsGroup)

	cats.GET("", r.catsHandler.GetAll)
	cats.GET("/:cat_id", r.catsHandler.GetById)
	cats.POST("", r.catsHandler.Create)
	cats.PATCH("/:cat_id", r.catsHandler.UpdateSalary)
	cats.DELETE("/:cat_id", r.catsHandler.Delete)

	missions.GET("", r.missionsHandler.GetMissions)
	missions.GET("/:mission_id", r.missionsHandler.GetMissionById)
	missions.POST("", r.missionsHandler.Create)
	missions.PATCH("/:mission_id/update", r.missionsHandler.UpdateStatus)
	missions.PATCH("/:mission_id/assign", r.missionsHandler.AssignCat)
	missions.DELETE("/:mission_id", r.missionsHandler.Delete)

	targets.POST("/:mission_id", r.targetsHandler.Create)
	targets.PATCH("/:target_id/update", r.targetsHandler.UpdateStatus)
	targets.PATCH("/:target_id/notes", r.targetsHandler.UpdateNote)
	targets.DELETE("/:target_id", r.targetsHandler.Delete)
}
