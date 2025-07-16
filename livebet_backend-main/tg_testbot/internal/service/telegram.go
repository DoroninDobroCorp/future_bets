package service

import (
	"context"
	"livebets/tg_testbot/cmd/config"
	"livebets/tg_testbot/internal/repository"
	"livebets/tg_testbot/internal/storage"
	"livebets/tg_testbot/internal/telegram"
	"livebets/tg_testbot/pkg/rdbms"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

type TelegramService struct {
	telegramTxStorage rdbms.TxStorage[repository.TelegramStorage]
	tgBot             *telegram.TelegramBot
	storage           *storage.Storage
	logger            *zerolog.Logger
}

func NewTelegramService(
	telegramTxStorage rdbms.TxStorage[repository.TelegramStorage],
	tgBot *telegram.TelegramBot,
	storage *storage.Storage,
	logger *zerolog.Logger,
) *TelegramService {
	return &TelegramService{
		telegramTxStorage: telegramTxStorage,
		tgBot:             tgBot,
		storage:           storage,
		logger:            logger,
	}
}

func (s *TelegramService) Run(ctx context.Context, cfg config.TGServiceConfig, wg *sync.WaitGroup) {
	defer wg.Done()

	seconds := cfg.Interval
	ticker := time.NewTicker(time.Duration(seconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.storage.NewFiles(ctx, s.telegramTxStorage.Storage(), s.tgBot, cfg.BetsPath)

		case <-ctx.Done():
			return
		}
	}
}
