package service

import (
	"fmt"
	"livebets/results/internal/api"
	"livebets/results/internal/entity"
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// TelegramService handles bot communication
type TelegramService struct {
	bot               *tgbotapi.BotAPI
	betService        *BetService
	adminChatID       int64
	adminUsername     string
	chatID            int64
	isWaitingForMatch bool
	batchBets         []entity.LogBetAccept
	batchDate         string
	currentBetKey     string
	currentOutcome    string
	isTestBet         bool // флаг для определения, что ожидаем ввод для тестовой ставки
}

func NewTelegramService(token string) (*TelegramService, error) {
	if token == "" {
		return nil, fmt.Errorf("telegram Bot API token is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating Telegram bot: %w", err)
	}

	// Set environment variables for admin chat ID and username
	adminChatIDStr := os.Getenv("TELEGRAM_ADMIN_CHAT_ID")
	adminUsername := os.Getenv("TELEGRAM_ADMIN_USERNAME")

	var adminChatID int64 = 0
	if adminChatIDStr != "" {
		adminChatID, err = strconv.ParseInt(adminChatIDStr, 10, 64)
		if err != nil {
			log.Printf("Warning: Cannot parse admin chat ID from environment: %v", err)
		}
	}

	service := &TelegramService{
		bot:           bot,
		adminChatID:   adminChatID,
		adminUsername: adminUsername,
	}

	return service, nil
}

func (s *TelegramService) SetBetService(betService *BetService) {
	s.betService = betService
}

func (s *TelegramService) SendMessage(text string) {
	if s.chatID == 0 {
		if s.adminChatID != 0 {
			s.chatID = s.adminChatID
		} else {
			log.Print("Warning: Cannot send message, no chat ID set")
			return
		}
	}

	msg := tgbotapi.NewMessage(s.chatID, text)
	_, err := s.bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message to Telegram: %v", err)
	}
}

func (s *TelegramService) Start() {
	if s.bot == nil {
		log.Print("Error: Telegram bot is not initialized")
		return
	}

	log.Printf("Starting Telegram bot: @%s", s.bot.Self.UserName)

	// Set up update channel
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := s.bot.GetUpdatesChan(u)

	// Process updates
	go func() {
		for update := range updates {
			if update.Message == nil {
				continue
			}

			// Store the chat ID so we can reply to the same chat
			s.chatID = update.Message.Chat.ID
			log.Printf("Получено сообщение от пользователя: %s, user ID: %d, chat ID: %d",
				update.Message.From.UserName,
				update.Message.From.ID,
				update.Message.Chat.ID)

			// Check if the message is from an admin user (if configured)
			if s.adminUsername != "" && update.Message.From.UserName != s.adminUsername {
				s.SendMessage("Sorry, you are not authorized to use this bot.")
				continue
			}

			// Process the message
			s.processMessage(update.Message)
		}
	}()
}

func (s *TelegramService) processMessage(message *tgbotapi.Message) {
	log.Printf("СУПЕР-ОТЛАДКА: Получено сообщение: '%s'", message.Text)
	log.Printf("СУПЕР-ОТЛАДКА: Состояние бота: isWaitingForMatch=%v, currentBetKey=%s, currentOutcome=%s",
		s.isWaitingForMatch, s.currentBetKey, s.currentOutcome)

	// Process commands
	if message.IsCommand() {
		log.Printf("СУПЕР-ОТЛАДКА: Получена команда: %s", message.Command())
		s.handleCommand(message.Command(), message.CommandArguments())
		return
	}

	// Default response for non-command messages
	log.Print("СУПЕР-ОТЛАДКА: Получено обычное сообщение (не в режиме ожидания, не команда)")
	s.SendMessage("Please use one of the available commands. Type /help for more information.")
}

func (s *TelegramService) handleCommand(command string, args string) {
	switch command {
	case "start":
		s.handleStartCommand()
	case "help":
		s.handleHelpCommand()
	case "check_main":
		s.handleCheckMainResultsCommand()
	case "check_test":
		s.handleCheckTestResultsCommand()
	case "recalc":
		s.handleRecalcCommand()
	case "fix54321":
		s.handleFixCommand()
	default:
		s.SendMessage("❌ Неизвестная команда. Используйте /help для вывода списка доступных команд.")
	}
}

func (s *TelegramService) handleStartCommand() {
	s.SendMessage("👋 Добро пожаловать в бот Results! Используйте /help для просмотра доступных команд.")
}

func (s *TelegramService) handleHelpCommand() {
	helpText := `📋 *Доступные команды:*

/check_main - запуск обработки результатов основных ставок
/check_test - запуск обработки результатов тестовых ставок
/recalc - обновить файл статистики по тестовым ставкам
`
	s.SendMessage(helpText)
}

func (s *TelegramService) handleFixCommand() {
	s.SendMessage("Fixing...")

	err := s.betService.Repo.FixTestDB()
	if err != nil {
		s.SendMessage(err.Error())
	} else {
		s.SendMessage("Fixed!")
	}
}

func (s *TelegramService) handleCheckMainResultsCommand() {
	if s.betService == nil {
		s.SendMessage("❌ Ошибка: Сервис ставок не инициализирован")
		return
	}

	// Обработка основной таблицы ставок
	err := s.betService.ProcessRecentBets()
	if err != nil {
		log.Printf("Error processing main bets: %v", err)
		s.SendMessage(fmt.Sprintf("Error processing main bets: %v", err))
	} else {
		s.SendMessage("Main bets processing completed successfully.")
		s.betService.PinnacleService.ClearCached()
	}
}

func (s *TelegramService) handleCheckTestResultsCommand() {
	if s.betService == nil {
		s.SendMessage("❌ Ошибка: Сервис ставок не инициализирован")
		return
	}

	// Обработка тестовой таблицы ставок
	err := s.betService.ProcessTestRecentBets()
	if err != nil {
		log.Printf("Error processing test bets: %v", err)
		s.SendMessage(fmt.Sprintf("Error processing test bets: %v", err))
	} else {
		s.SendMessage("Test bets processing completed successfully.")
		s.betService.PinnacleService.ClearCached()
	}
}

func (s *TelegramService) handleRecalcCommand() {
	rows, err := s.betService.Repo.GetTestBets()
	if err != nil {
		s.SendMessage(fmt.Sprintf("Error getting recalc bets: %v", err))
		return
	}

	s.SendMessage(fmt.Sprintf("Got %d test bets for statistic", len(rows)))

	table := api.ProcessTable(rows)

	err = api.SaveToCSV(table)
	if err != nil {
		s.SendMessage(fmt.Sprintf("Error saving recalc bets: %v", err))
	}

	s.SendMessage("File statistic.csv updated successfully.")
}
