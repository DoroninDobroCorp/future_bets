package main

import (
	"context"
	"github.com/rs/zerolog"
	"livebets/parse_maxbet/cmd/config"
	"livebets/parse_maxbet/internal/api"
	"livebets/parse_maxbet/internal/sender"
	"livebets/parse_maxbet/internal/service"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Parse_Maxbet")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	api := api.New(appConfig.APIConfig)

	sender := sender.New(appConfig.SenderConfig)

	service := service.New(api, sender, &logger)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go service.RunLive(ctx, appConfig.APIConfig, wg)

	http.HandleFunc("/health", HealthCheckHandler)
	server := &http.Server{Addr: ":" + appConfig.Port}

	go func() {
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancelFunc()
	wg.Wait()

	if err = server.Shutdown(context.Background()); err != nil {
		logger.Fatal().Err(err).Msg("failed to stop server")
	}

	logger.Info().Msg(">> Stopping Parse_Maxbet")
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
