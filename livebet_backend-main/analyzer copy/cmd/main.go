package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"livebets/analazer/cmd/config"
	"livebets/analazer/internal/entity"
	"livebets/analazer/internal/handler"
	priceStorage "livebets/analazer/internal/price-storage"
	"livebets/analazer/internal/receiver"
	"livebets/analazer/internal/repository"
	"livebets/analazer/internal/sender"
	"livebets/analazer/internal/service"
	"livebets/analazer/migrations"
	"livebets/analazer/pkg/pgsql"
	"livebets/analazer/pkg/rdbms"
	redis_client "livebets/analazer/pkg/redis"
	"livebets/analazer/pkg/server"

	"github.com/rs/zerolog"
)

// Главная функция
func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting Analyzer")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	// Connect to postgres
	postgres, err := pgsql.New(appConfig.PostgresConfig.ConnectionString())
	if err != nil {
		cancelFunc()
		logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	logger.Info().Msg("Connected to Postgres")
	defer postgres.Close()

	if err = ensureMigrate(ctx, appConfig.PostgresConfig.ConnectionString(), &logger); err != nil {
		cancelFunc()
		logger.Fatal().Err(err).Msg("unable to migrate service database")
	}

	// Connect to redis
	redisClient, err := redis_client.NewRedis(ctx, appConfig.RedisConfig)
	if err != nil {
		cancelFunc()
		logger.Fatal().Err(err).Msg("failed to connect to redis")
	}
	logger.Info().Msg("Connected to Redis")
	defer redisClient.Close()

	wg := &sync.WaitGroup{}

	// Initialization
	pairsMatchTxStorage := rdbms.NewPgTxStorage(postgres.Pool, repository.NewPairsMatchingPGStorage)
	priceStorage := priceStorage.NewPriceStorage()
	wg.Add(1)
	go priceStorage.CleanByTimeout(ctx, appConfig.PriceStorage, wg)

	receiveChan := make(chan entity.ReceivedMsg, 100)
	defer close(receiveChan)
	sendChan := make(chan []entity.ResponsePair, 500)
	defer close(sendChan)

	receiverPinnacle := receiver.NewReceiver(receiveChan)
	receiverOthers := receiver.NewReceiver(receiveChan)

	// Запускаем два сервера: один для Pinnacle, другой для всех остальных
	go receiverPinnacle.StartParserServer(appConfig.Port.Pinnacle)
	go receiverOthers.StartParserServer(appConfig.Port.Other)

	pairsMatchService := service.NewPairsMatchingService(pairsMatchTxStorage, redisClient, receiveChan, sendChan, priceStorage, &logger)
	wg.Add(1)
	go pairsMatchService.Run(ctx, appConfig.PairsMatching, wg)

	// TODO: rewrite start server
	sender := sender.NewSender(sendChan, &logger)
	go sender.StartServer(appConfig.Port.Sender)
	wg.Add(1)
	go sender.RunHub(ctx, wg)

	priceService := service.NewPriceService(priceStorage, &logger)

	handlers := handler.NewHandler(priceService, pairsMatchService)

	srv := new(server.Server)
	go func() {
		logger.Info().Msgf("starting server on port = %d", appConfig.Port.Server)
		if err := srv.Run(fmt.Sprintf("%d", appConfig.Port.Server), handlers.InitRoutes()); err != nil {
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

	logger.Info().Msg(">> Stopping Analyzer")
}

func ensureMigrate(ctx context.Context, connString string, logger *zerolog.Logger) error {
	migrator := pgsql.Migrator{
		Context:          ctx,
		Logger:           logger,
		SchemaName:       "analyzer",
		TableName:        "migrations",
		MigrationsFS:     migrations.MigrationsFS,
		ConnectionString: connString,
		RootFS:           ".",
		Driver:           "postgres",
	}

	return migrator.Run()
}
