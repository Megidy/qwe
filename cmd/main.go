package main

import (
	"github.com/Megidy/cats/cmd/app"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msgf("Hello world")

	app, err := app.NewApp()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create application instance")
	}
	app.Run()

}
