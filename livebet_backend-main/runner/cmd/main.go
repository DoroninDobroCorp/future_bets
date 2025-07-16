package main

import (
	"context"
	"livebets/runner/cmd/config"
	"livebets/runner/internal/api"
	"livebets/runner/internal/handler"
	"livebets/runner/internal/service"
	"livebets/runner/internal/storage"
	"livebets/runner/pkg/server"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Runner")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	wg := &sync.WaitGroup{}

	bookmakersStorage := storage.NewBookmakerStorage(appConfig.Bookmakers)
	statusStorage := storage.NewStatusStorage()
	parserAPI := api.NewParserAPI()

	commandService := service.NewCommandService(appConfig.CommandConfig, bookmakersStorage, &logger)
	wg.Add(1)
	go commandService.Run(ctx, wg)

	statusService := service.NewStatusService(appConfig.StatusConfig, parserAPI, bookmakersStorage, statusStorage, &logger)
	wg.Add(1)
	go statusService.Run(ctx, wg)

	handlers := handler.NewHandler(commandService, statusService)

	srv := new(server.Server)
	go func() {
		logger.Info().Msgf("starting server on port = %s", appConfig.CommandConfig.Port)
		if err := srv.Run(appConfig.CommandConfig.Port, handlers.InitRoutes()); err != nil {
			logger.Error().Err(err).Msg("error occured while running http server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancelFunc()
	wg.Wait()

	if err = srv.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("error occured on server shutting down")
	}
	logger.Info().Msg(">> Runner")
}
