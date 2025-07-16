package main

import (
	"context"
	"livebets/calculator/cmd/config"
	"livebets/calculator/internal/api"
	"livebets/calculator/internal/handler"
	"livebets/calculator/internal/repository"
	"livebets/calculator/internal/service"
	"livebets/calculator/migrations"
	"livebets/calculator/pkg/pgsql"
	"livebets/calculator/pkg/rdbms"
	"livebets/calculator/pkg/server"
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
	logger.Info().Msg(">> Starting Calculator")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	// Connect to postgres
	postgres, err := pgsql.New(appConfig.PostgresConfig.ConnectionString())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to postgres")
	}
	logger.Info().Msg("Connected to Postgres")
	defer postgres.Close()

	if err = ensureMigrate(ctx, appConfig.PostgresConfig.ConnectionString(), &logger); err != nil {
		cancelFunc()
		logger.Fatal().Err(err).Msg("unable to migrate service database")
	}

	wg := &sync.WaitGroup{}

	analizerAPI := api.NewAnalizerAPI(appConfig.AnalyzerAPI)
	analizerPrematchAPI := api.NewAnalizerPrematchAPI(appConfig.AnalyzerPrematchAPI)
	logsTxStorage := rdbms.NewPgTxStorage(postgres.Pool, repository.NewLogsPGStorage)
	logsService := service.NewLogsService(logsTxStorage, analizerAPI, analizerPrematchAPI, &logger)
	if err = logsService.InitializeTotalBetPercents(ctx); err != nil {
		logger.Error().Err(err).Msg("unable to read database for calc percent")
	}

	wg.Add(1)
	go logsService.CleanCaches(ctx, appConfig.LogsService, wg)

	handlers := handler.NewHandler(logsService)

	srv := new(server.Server)
	go func() {
		logger.Info().Msgf("starting server on port = %s", "7010")
		if err := srv.Run("7010", handlers.InitRoutes()); err != nil {
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
	logger.Info().Msg(">> Stopping Calculator")
}

func ensureMigrate(ctx context.Context, connString string, logger *zerolog.Logger) error {
	migrator := pgsql.Migrator{
		Context:          ctx,
		Logger:           logger,
		SchemaName:       "calculator",
		TableName:        "migrations",
		MigrationsFS:     migrations.MigrationsFS,
		ConnectionString: connString,
		RootFS:           ".",
		Driver:           "postgres",
	}

	return migrator.Run()
}
