package telegram

import (
	"context"
	"fmt"
	"livebets/tg_testbot/cmd/config"
	"livebets/tg_testbot/internal/repository"
	"livebets/tg_testbot/pkg/rdbms"
	"strconv"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

type TelegramBot struct {
	bot               *tgbotapi.BotAPI
	telegramTxStorage rdbms.TxStorage[repository.TelegramStorage]
	logger            *zerolog.Logger
	subscribers       map[int64]string
	groupsByBookmaker map[string]int64
}

func NewTelegramBot(
	bot *tgbotapi.BotAPI,
	telegramTxStorage rdbms.TxStorage[repository.TelegramStorage],
	logger *zerolog.Logger,
	groupIds config.GroupsByBookmaker,
) *TelegramBot {
	groupsByBookmaker := map[string]int64{
		"Lobbet":     groupIds.LobbetGroup,
		"Ladbrokes":  groupIds.LadbrokesGroup,
		"Ladbrokes2": groupIds.Ladbrokes2Group,
		"Unibet":     groupIds.UnibetGroup,
		"Starcasino": groupIds.StarcasinoGroup,
	}

	return &TelegramBot{
		bot:               bot,
		telegramTxStorage: telegramTxStorage,
		logger:            logger,
		subscribers:       make(map[int64]string),
		groupsByBookmaker: groupsByBookmaker,
	}
}

func (t *TelegramBot) SendFile(fileName string, file []byte) error {
	splited := strings.Split(fileName, "_")
	if len(splited) != 7 {
		return fmt.Errorf("invalid file name format")
	}

	bkName := splited[3]
	roi, _ := strconv.Atoi(strings.ReplaceAll(splited[6], ".csv", ""))

	fmt.Printf("Splitted: %v\nbkName: '%s' | roi: '%s'\n", splited, bkName, roi)

	fileBytes := tgbotapi.FileBytes{Name: fileName, Bytes: file}

	fileWasSent := false

	if roi < 3 {
		for chatID, userName := range t.subscribers {
			userStr := userString(userName, chatID)

			newDoc := tgbotapi.NewDocument(chatID, fileBytes)

			_, err := t.bot.Send(newDoc)
			if err != nil {
				t.logger.Error().Err(err).Msg("telegram bot failed to send file. " + userStr + " File: " + fileName)
				continue
			}

			t.logger.Info().Msg("File was sent to user. " + userStr + " File: " + fileName)
			fileWasSent = true
		}
	} else {
		bookmakerGroupId, ok := t.groupsByBookmaker[bkName]

		if ok {
			userStr := userString(bkName, bookmakerGroupId)

			newDoc := tgbotapi.NewDocument(bookmakerGroupId, fileBytes)

			_, err := t.bot.Send(newDoc)
			if err != nil {
				t.logger.Error().Err(err).Msg("telegram bot failed to send file. " + userStr + " File: " + fileName)
				return err
			}

			t.logger.Info().Msg("File was sent to user. " + userStr + " File: " + fileName)
			fileWasSent = true

		} else {
			for chatID, userName := range t.subscribers {
				userStr := userString(userName, chatID)

				newDoc := tgbotapi.NewDocument(chatID, fileBytes)

				_, err := t.bot.Send(newDoc)
				if err != nil {
					t.logger.Error().Err(err).Msg("telegram bot failed to send file. " + userStr + " File: " + fileName)
					continue
				}

				t.logger.Info().Msg("File was sent to user. " + userStr + " File: " + fileName)
				fileWasSent = true
			}
		}
	}

	if !fileWasSent {
		return fmt.Errorf("there are no users to send file")
	}
	return nil
}

func (t *TelegramBot) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	t.subscribers, _ = t.GetAllSubscribers(ctx)
	t.logger.Info().Msg(fmt.Sprintf("Bot has %d subscribers.", len(t.subscribers)))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updatesChannel := t.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case update := <-updatesChannel:
			t.parseMessage(ctx, update.Message)

		case <-ctx.Done():
			return
		}
	}
}

func (t *TelegramBot) parseMessage(ctx context.Context, message *tgbotapi.Message) {
	if message == nil {
		return
	}

	username := "@" + message.Chat.UserName
	chatID := message.Chat.ID

	if message.Chat.Type == "private" {
		if strings.HasSuffix(message.Text, "nsubscribe") {
			t.unsubscribeUser(ctx, username, chatID)
		} else {
			t.subscribeUser(ctx, username, chatID)
		}
	}
}

func (t *TelegramBot) subscribeUser(ctx context.Context, userName string, chatID int64) {
	userStr := userString(userName, chatID)

	_, exists := t.subscribers[chatID]

	text := ""
	if exists {
		text = "You have already subscribed to our mailing list of CSV-files.\nTo unsubscribe, type \"Unsubscribe\"."
	} else {
		err := t.telegramTxStorage.Storage().AddSubscriber(ctx, userName, chatID)
		if err != nil {
			t.logger.Error().Err(err).Msg("filed to add subscriber. " + userStr)
			return
		}

		t.subscribers[chatID] = userName

		t.logger.Info().Msg("User subscribed. " + userStr)

		text = "You have subscribed to our mailing list of CSV-files.\n" + userStr
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error().Err(err).Msg("telegram bot failed to send subscribe message. " + userStr)
	}
}

func (t *TelegramBot) unsubscribeUser(ctx context.Context, userName string, chatID int64) {
	userStr := userString(userName, chatID)

	_, exists := t.subscribers[chatID]

	text := ""
	if exists {
		err := t.telegramTxStorage.Storage().DeleteSubscriber(ctx, chatID)
		if err != nil {
			t.logger.Error().Err(err).Msg("filed to delete subscriber. " + userStr)
			return
		}

		delete(t.subscribers, chatID)

		t.logger.Info().Msg("User unsubscribed. " + userStr)

		text = "You have unsubscribed from our mailing list for CSV-files." + "\n" + userStr
	} else {
		text = "You have already unsubscribed from our mailing list for CSV-files."
	}

	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := t.bot.Send(msg); err != nil {
		t.logger.Error().Err(err).Msg("telegram bot failed to send unsubscribe message. " + userStr)
	}
}

func (t *TelegramBot) GetAllSubscribers(ctx context.Context) (map[int64]string, error) {
	return t.telegramTxStorage.Storage().GetAllSubscribers(ctx)
}

func userString(userName string, chatID int64) string {
	return fmt.Sprintf("Username: %s ChatID: %d", userName, chatID)
}
