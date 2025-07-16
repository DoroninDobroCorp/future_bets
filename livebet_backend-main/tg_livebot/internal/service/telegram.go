package service

import (
	"context"
	"livebets/tg_livebot/cmd/config"
	"livebets/tg_livebot/internal/repository"
	"livebets/tg_livebot/internal/storage"
	"livebets/tg_livebot/internal/telegram"
	"livebets/tg_livebot/pkg/rdbms"
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
