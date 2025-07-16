package main

import (
	"context"
	"livebets/parse_starcasino/cmd/config"
	"livebets/parse_starcasino/internal/api"
	"livebets/parse_starcasino/internal/sender"
	"livebets/parse_starcasino/internal/service"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// Основная функция
func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Parse_StarCasino")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	api := api.New(appConfig.APIConfig)

	sender := sender.New(appConfig.SenderConfig)

	service := service.New(api, sender, &logger)

	wg := &sync.WaitGroup{}

	if appConfig.ParseLive {
		wg.Add(1)
		go service.RunLive(ctx, appConfig.APIConfig, wg)
	} else {
		wg.Add(1)
		go service.RunPreMatch(ctx, appConfig.APIConfig, wg)
	}

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

	logger.Info().Msg(">> Stopping Parse_StarCasino")
}
