package repository

import (
	"context"
	"fmt"
	"livebets/tg_testbot/pkg/rdbms"
	"time"
)

type TelegramStorage interface {
	// Files
	FindFile(ctx context.Context, fileName string) (bool, error)
	SaveNewFile(ctx context.Context, fileName string, fileTime time.Time) error
	// Subscribers
	AddSubscriber(ctx context.Context, userName string, chatID int64) error
	DeleteSubscriber(ctx context.Context, chatID int64) error
	GetAllSubscribers(ctx context.Context) (map[int64]string, error)
}

type TelegramPGStorage struct {
	handler rdbms.Executor
}

func NewTelegramPGStorage(handler rdbms.Executor) TelegramStorage {
	return &TelegramPGStorage{
		handler: handler,
	}
}
func (t *TelegramPGStorage) FindFile(ctx context.Context, fileName string) (bool, error) {

	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE filename = $1)", botFilesTable)

	var exists bool
	err := t.handler.QueryRow(ctx, query, fileName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (t *TelegramPGStorage) SaveNewFile(ctx context.Context, fileName string, fileTime time.Time) error {
	query := fmt.Sprintf("INSERT INTO %s (filename, date) VALUES ($1, $2)", botFilesTable)

	_, err := t.handler.Exec(ctx, query, fileName, fileTime)
	if err != nil {
		return err
	}

	return nil
}

func (t *TelegramPGStorage) AddSubscriber(ctx context.Context, userName string, chatID int64) error {
	query := fmt.Sprintf("INSERT INTO %s (username, chat_id) VALUES ($1, $2)", botSubscribersTable)

	_, err := t.handler.Exec(ctx, query, userName, chatID)

	return err
}

func (t *TelegramPGStorage) DeleteSubscriber(ctx context.Context, chatID int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE chat_id = $1", botSubscribersTable)

	_, err := t.handler.Exec(ctx, query, chatID)

	return err
}

func (t *TelegramPGStorage) GetAllSubscribers(ctx context.Context) (map[int64]string, error) {
	query := fmt.Sprintf("SELECT username, chat_id FROM %s", botSubscribersTable)

	rows, err := t.handler.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[int64]string)

	for rows.Next() {
		var username string
		var chatID int64

		err = rows.Scan(&username, &chatID)
		if err != nil {
			return nil, err
		}

		users[chatID] = username
	}

	return users, nil
}
