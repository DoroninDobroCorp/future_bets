package main

import (
	"context"
	"livebets/tg_livebot/cmd/config"
	"livebets/tg_livebot/internal/repository"
	"livebets/tg_livebot/internal/service"
	"livebets/tg_livebot/internal/storage"
	"livebets/tg_livebot/internal/telegram"
	"livebets/tg_livebot/migrations"
	"livebets/tg_livebot/pkg/pgsql"
	"livebets/tg_livebot/pkg/rdbms"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Init logger and config
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger.Info().Msg(">> Starting TG_livebot")
	appConfig, err := config.ProvideAppMPConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load app configuration")
	}

	// Connect to database
	postgres, err := pgsql.New(appConfig.PostgresConfig.ConnectionString())
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	logger.Info().Msg("Connected to Postgres")
	defer postgres.Close()

	if err = ensureMigrate(ctx, appConfig.PostgresConfig.ConnectionString(), &logger); err != nil {
		cancelFunc()
		logger.Fatal().Err(err).Msg("unable to migrate service database")
	}

	// Init telegram bot
	bot, err := tgbotapi.NewBotAPI(appConfig.TelegramBotConfig.Token)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect telegram API")
	}
	logger.Info().Msg("Connected to telegram API. Bot account : " + bot.Self.UserName)

	telegramTxStorage := rdbms.NewPgTxStorage(postgres.Pool, repository.NewTelegramPGStorage)

	storage := storage.NewStorage(&logger)

	telegramBot := telegram.NewTelegramBot(bot, telegramTxStorage, &logger)

	telegramService := service.NewTelegramService(telegramTxStorage, telegramBot, storage, &logger)

	wg := &sync.WaitGroup{}

	// Обработка файлов
	wg.Add(1)
	go telegramService.Run(ctx, appConfig.TGServiceConfig, wg)

	// Обработка списка рассылки
	wg.Add(1)
	go telegramBot.Run(ctx, wg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	cancelFunc()
	wg.Wait()

	logger.Info().Msg(">> Stopping TG_livebot")
}

func ensureMigrate(ctx context.Context, connString string, logger *zerolog.Logger) error {
	migrator := pgsql.Migrator{
		Context:          ctx,
		Logger:           logger,
		SchemaName:       "tg_livebot",
		TableName:        "migrations",
		MigrationsFS:     migrations.MigrationsFS,
		ConnectionString: connString,
		RootFS:           ".",
		Driver:           "postgres",
	}

	return migrator.Run()
}
