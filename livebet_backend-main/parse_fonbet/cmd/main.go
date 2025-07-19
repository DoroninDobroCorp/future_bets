package main

import (
	"context"
	"livebets/parse_fonbet/cmd/config"
	"livebets/parse_fonbet/internal/api"
	"livebets/parse_fonbet/internal/entity"
	"livebets/parse_fonbet/internal/sender"
	"livebets/parse_fonbet/internal/service"
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

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Parse_Fonbet")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	sendChan := make(chan entity.ResponseGame, 50)
	defer close(sendChan)

	api := api.NewFonbetAPI(appConfig.FonbetConfig)
	sender := sender.NewSender(appConfig.SenderConfig, sendChan)
	service := service.NewGeneralService(api, sendChan, &logger)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go sender.SendingToAnalyzer(ctx, wg)
	wg.Add(1)
	go service.RunFootball(ctx, appConfig.FonbetConfig, wg, true)
	wg.Add(1)
	go service.RunTennis(ctx, appConfig.FonbetConfig, wg, true)

	http.HandleFunc("/health", HealthCheckHandler)
	if err := http.ListenAndServe(":"+appConfig.Port, nil); err != nil {
		logger.Fatal().Err(err).Msg("failed to start server")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancelFunc()
	wg.Wait()

	logger.Info().Msg(">> Stopping Parse_Fonbet")
}
