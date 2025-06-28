package app

import (
	"context"
	"fmt"

	"github.com/Megidy/cats/internal/config"
	"github.com/Megidy/cats/internal/dbconnection"
	"github.com/Megidy/cats/internal/infrastructure/catcaller"
	"github.com/Megidy/cats/internal/infrastructure/repository"
	"github.com/Megidy/cats/internal/interface/http"
	"github.com/Megidy/cats/internal/interface/service"
	httpserver "github.com/Megidy/cats/pkg/servers/http"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App interface {
	Run() error
	Shutdown()
}

type app struct {
	httpServer *httpserver.HttpServer
	pool       *pgxpool.Pool
}

func NewApp() (App, error) {

	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	pool, err := dbconnection.NewPostgreSQLConnectionPool(context.Background(), cfg.PostgresURI)
	if err != nil {
		return nil, err
	}

	server := httpserver.NewHttpServer(cfg.HttpServerPort)

	catCaller := catcaller.NewCatCaller()

	catsRepo := repository.NewCatsRepository(pool)
	missionsRepo := repository.NewMissionsRepository(pool)
	targetsRepo := repository.NewTargetsRepository(pool)

	catsService := service.NewCatService(catsRepo, catCaller)
	missionsService := service.NewMissionService(missionsRepo, catsRepo)
	targetService := service.NewTargetService(targetsRepo, missionsRepo)

	catsHandler := http.NewCatsHandler(catsService)
	missionsHandler := http.NewMissionsHandler(missionsService)
	targetsHandler := http.NewTargetsHandler(targetService)

	router := http.NewRouter(
		server.Echo,
		catsHandler,
		missionsHandler,
		targetsHandler,
	)
	router.RegisterRoutes()

	return &app{
		httpServer: server,
		pool:       pool,
	}, nil
}

func (a *app) Run() error {
	return a.httpServer.Run()

}
func (a *app) Shutdown() {}
