package storage

import (
	"context"
	"github.com/rs/zerolog"
	"livebets/tg_livebot/internal/repository"
	"livebets/tg_livebot/internal/telegram"
	"os"
	"path"
)

type Storage struct {
	logger *zerolog.Logger
}

func NewStorage(
	logger *zerolog.Logger,
) *Storage {
	return &Storage{
		logger: logger,
	}
}

// Working with files
func (s *Storage) NewFiles(ctx context.Context, telegramStorage repository.TelegramStorage, tgBot *telegram.TelegramBot, filePath string) {

	files, err := os.ReadDir(filePath)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read CSV-files directory: " + filePath)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		fileInfo, _ := file.Info()
		fileTime := fileInfo.ModTime()

		// Проверяем наличие файла в таблице bot_files
		fileWasFound, err := telegramStorage.FindFile(ctx, fileName)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to find file. File: " + fileName)
			continue
		}

		if !fileWasFound {
			// Файла нет в базе. Отправляем файл в Telegram
			if err := s.uploadFile(tgBot, filePath, fileName); err != nil {
				s.logger.Error().Err(err).Msg("failed to upload file. File: " + fileName)
				// Файл не удалось отправить. Оставляем файл на диске
				continue
			}

			// Записываем файл в таблицу
			err = telegramStorage.SaveNewFile(ctx, fileName, fileTime)
			if err != nil {
				s.logger.Error().Err(err).Msg("failed to save new file. File: " + fileName)
				continue
			}
			s.logger.Info().Msg("File was saved to database. File: " + fileName)
		}

		err = s.deleteFile(filePath, fileName)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to delete file : " + fileName)
			continue
		}
		s.logger.Info().Msg("File was deleted. File: " + fileName)
	}
}

func (s *Storage) deleteFile(filePath, fileName string) error {
	fullFileName := path.Join(filePath, fileName)

	err := os.Remove(fullFileName)

	return err
}

func (s *Storage) uploadFile(tgBot *telegram.TelegramBot, filePath, fileName string) error {
	fullFileName := path.Join(filePath, fileName)

	file, err := os.ReadFile(fullFileName)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to read file : " + fileName)
		return err
	}

	err = tgBot.SendFile(fileName, file)

	return err
}
