package main

import (
	"fmt"
	"livebets/results/internal/api"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"livebets/results/internal/repository"
	"livebets/results/internal/service"
)

func main() {
	// Проверяем наличие аргументов командной строки
	if len(os.Args) > 1 {
		// Получаем первый аргумент
		command := os.Args[1]

		// Здесь могут быть обработаны другие команды в будущем
		log.Printf("Command argument provided but not recognized: %s", command)
	}

	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	log.Printf("Working directory: %s", wd)

	// Load .env file
	envPath := filepath.Join(wd, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	} else {
		log.Printf("Loaded .env from %s", envPath)
	}

	// Database connection
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USERNAME")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	pgClient, err := repository.NewPostgresClient(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("DB connection successful")

	// Telegram service
	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if telegramToken == "" {
		telegramToken = "default_token"
		log.Println("Warning: TELEGRAM_BOT_TOKEN not set")
	}

	telegramService, err := service.NewTelegramService(telegramToken)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram service: %v", err)
	}

	pinnacleService := service.NewPinnacleService(
		"AG1677099",
		"Df2321sT43",
		true,
		"socks5://AllivanService:ProxyAll1van@154.7.188.74:5089",
	)

	// Set up Bet service
	betService := service.NewBetService(pgClient, pinnacleService)
	betService.TelegramService = telegramService
	pinnacleService.SetTelegramLogger(telegramService)
	telegramService.SetBetService(betService)

	// Start the Telegram bot
	telegramService.Start()

	// Send starting message
	telegramService.SendMessage("Service started")

	// Main loop
	log.Println("Service started")
	for {
		now := time.Now()
		if now.Hour() == 7 && now.Minute() == 20 {
			log.Println("Processing recent bets...")
			telegramService.SendMessage("Starting daily bet processing...")

			// Обработка основной таблицы ставок
			err := betService.ProcessRecentBets()
			if err != nil {
				log.Printf("Error processing main bets: %v", err)
				telegramService.SendMessage(fmt.Sprintf("Error processing main bets: %v", err))
			} else {
				telegramService.SendMessage("Main bets processing completed successfully.")
			}

			// Обработка тестовой таблицы ставок
			err = betService.ProcessTestRecentBets()
			if err != nil {
				log.Printf("Error processing test bets: %v", err)
				telegramService.SendMessage(fmt.Sprintf("Error processing test bets: %v", err))
			} else {
				telegramService.SendMessage("Test bets processing completed successfully.")
			}

			telegramService.SendMessage("Daily bet processing completed for both tables.")
			pinnacleService.ClearCached()

			// Обновление таблицы статистики
			rows, err := betService.Repo.GetTestBets()
			if err != nil {
				telegramService.SendMessage(fmt.Sprintf("Error getting recalc bets: %v", err))
				return
			}

			telegramService.SendMessage(fmt.Sprintf("Got %d test bets for statistic", len(rows)))

			table := api.ProcessTable(rows)

			err = api.SaveToCSV(table)
			if err != nil {
				telegramService.SendMessage(fmt.Sprintf("Error saving recalc bets: %v", err))
			}

			time.Sleep(61 * time.Second)
		}

		time.Sleep(60 * time.Second)
	}
}
