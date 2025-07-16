package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"livebets/results/internal/calculator"
	"livebets/results/internal/entity"
	"livebets/results/internal/repository"
)

type BetService struct {
	Repo            *repository.PostgresClient
	PinnacleService *PinnacleService
	TelegramService *TelegramService
	LogToTelegram   bool
}

func NewBetService(repo *repository.PostgresClient, pinnacleService *PinnacleService) *BetService {
	return &BetService{
		Repo:            repo,
		PinnacleService: pinnacleService,
		LogToTelegram:   true,
	}
}

type BookmakerSummary struct {
	TotalBets   int
	TotalSum    float64
	TotalProfit float64
}

func printBetSummary(bet *entity.LogBetAccept, outcomeStr string, betSum, coef float64) {
	pair, ok := bet.Data["pair"].(map[string]interface{})
	if !ok {
		log.Printf("DEBUG: Failed to extract 'pair' field from bet.Data for bet %s", bet.KeyMatch)
		return
	}
	first, ok := pair["first"].(map[string]interface{})
	if !ok {
		log.Printf("DEBUG: Failed to extract 'first' field from pair for bet %s", bet.KeyMatch)
		return
	}
	// Get team names from first (homeName and awayName)
	homeName, _ := first["homeName"].(string)
	awayName, _ := first["awayName"].(string)
	log.Printf("DEBUG: Teams: %s - %s, Sum: %.2f, Outcome: %s, Coefficient: %.2f",
		homeName, awayName, betSum, outcomeStr, coef)
}

func (s *BetService) ProcessBet(bet *entity.LogBetAccept) (calculator.OutcomeResult, error) {
	logMsg := fmt.Sprintf("Starting to process bet with key_match: %s", bet.KeyMatch)
	log.Print(logMsg)

	// Логируем начало обработки ставки
	log.Printf("DEBUG: Начало обработки ставки %s", bet.KeyMatch)

	// Extract outcome from bet.Data (expected structure: bet.Data -> "pair" -> "outcome" -> "outcome")
	var outcomeStr string
	var sport string
	if pair, ok := bet.Data["pair"].(map[string]interface{}); ok {
		sport = pair["sportName"].(string)
		if outcomeObj, ok := pair["outcome"].(map[string]interface{}); ok {
			if oStr, ok := outcomeObj["outcome"].(string); ok && oStr != "" {
				outcomeStr = oStr
			} else {
				errMsg := "outcome field inside pair.outcome not found or empty"
				log.Printf("Error: %s", errMsg)
				return calculator.OutcomeResult{Status: "error"}, errors.New(errMsg)
			}
		} else {
			errMsg := "outcome field in pair not found"
			log.Printf("Error: %s", errMsg)
			return calculator.OutcomeResult{Status: "error"}, errors.New(errMsg)
		}
	} else {
		errMsg := "pair field not found"
		log.Printf("Error: %s", errMsg)
		return calculator.OutcomeResult{Status: "error"}, errors.New(errMsg)
	}
	debugMsg := fmt.Sprintf("DEBUG: Extracted outcome for bet %s: %s", bet.KeyMatch, outcomeStr)
	log.Print(debugMsg)

	// Extract sum and coefficient
	sumVal, ok := bet.Data["sum"].(float64)
	if !ok {
		errMsg := "bet sum not found or has invalid type"
		return calculator.OutcomeResult{Status: "error"}, errors.New(errMsg)
	}
	coefVal, ok := bet.Data["coef"].(float64)
	if !ok {
		errMsg := "coefficient not found or has invalid type"
		log.Printf("Error: %s", errMsg)
		return calculator.OutcomeResult{Status: "error"}, errors.New(errMsg)
	}
	printBetSummary(bet, outcomeStr, sumVal, coefVal)

	// Get fixture through PinnacleService (with cache inside CallFixtureSettled)
	fixture, err := s.PinnacleService.CallFixtureSettled(bet)
	if err != nil {
		log.Printf("Failed to get fixture: %v. Status pending.", err)

		// Подробно логируем доступные данные ставки для диагностики
		betData, _ := json.Marshal(bet.Data)
		log.Printf("DEBUG: Данные ставки для диагностики: %s", string(betData))

		// Проверяем, есть ли у нас хотя бы базовую информацию о матче
		var matchInfo string
		var isLive bool
		if pair, ok := bet.Data["pair"].(map[string]interface{}); ok {
			if isLiveVal, ok := pair["isLive"].(bool); ok {
				isLive = isLiveVal
			}
			if first, ok := pair["first"].(map[string]interface{}); ok {
				home, _ := first["homeName"].(string)
				away, _ := first["awayName"].(string)
				matchInfo = fmt.Sprintf("%s vs %s", home, away)
			}
		}

		if !isLive && bet.RealProfit == nil {
			newTime := time.Date(
				time.Now().Year(),
				time.Now().Month(),
				time.Now().Day(),
				bet.CreatedAt.Hour(),
				bet.CreatedAt.Minute(),
				bet.CreatedAt.Second(),
				bet.CreatedAt.Nanosecond(),
				bet.CreatedAt.Location(),
			)

			err = s.Repo.UpdateBetTime(bet.KeyOutcome, newTime)
			if err != nil {
				log.Printf("Failed to update bet time: %v", err)
			}
		}

		if matchInfo != "" {
			log.Printf("DEBUG: Информация о матче: %s", matchInfo)
		}
		return calculator.OutcomeResult{Status: "pending"}, nil
	}

	// Логируем информацию о найденном матче
	log.Printf("DEBUG: Получены данные матча от Pinnacle для ставки %s. ID: %d, Количество периодов: %d",
		bet.KeyMatch, fixture.ID, len(fixture.Periods))

	for _, period := range fixture.Periods {
		log.Printf("DEBUG: Период %d: счет %d:%d", period.Number, period.Team1Score, period.Team2Score)
	}

	debugMsg = fmt.Sprintf("DEBUG: Got fixture for bet %s", bet.KeyMatch)
	log.Print(debugMsg)

	// Calculate bet outcome
	result := calculator.GetOutcomeResult(bet, fixture, outcomeStr, sumVal, coefVal, sport)
	return result, nil
}

func (s *BetService) ProcessRecentBets() error {
	logMsg := fmt.Sprintf("Processing bets for the last day")
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	// Get bets for the specified period
	bets, err := s.Repo.GetYesterdayBets()
	if err != nil {
		errMsg := fmt.Sprintf("Error getting recent bets: %v", err)
		log.Print(errMsg)
		if s.LogToTelegram && s.TelegramService != nil {
			s.TelegramService.SendMessage(errMsg)
		}
		return err
	}

	logMsg = fmt.Sprintf("Found %d bets for the last day", len(bets))
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	// Process each bet
	totalBets := len(bets)
	processedBets := 0
	for _, bet := range bets {
		result, err := s.ProcessBet(bet)
		if err != nil {
			log.Printf("Error processing bet %s: %v", bet.KeyOutcome, err)
		} else {
			log.Printf("Processed bet %s\nOutcome: `%s`\nCoef: %.2f\nResult: %s", bet.KeyOutcome, bet.CorrectData["outcome"], bet.Data["coef"], result.Status)
			if s.LogToTelegram && s.TelegramService != nil {
				s.TelegramService.SendMessage(fmt.Sprintf("Processed bet %s\nOutcome: `%s`\nCoef: %.2f\nResult: %s", bet.KeyOutcome, bet.CorrectData["outcome"], bet.Data["coef"], result.Status))
			}

			// Извлекаем сумму ставки, коэффициент и ROI
			sumVal := bet.Data["sum"].(float64)
			coefVal := bet.Data["coef"].(float64)
			var roi float64
			if pair, ok := bet.Data["pair"].(map[string]interface{}); ok {
				if outcomeObj, ok := pair["outcome"].(map[string]interface{}); ok {
					if roiVal, ok := outcomeObj["roi"].(float64); ok {
						roi = roiVal
					}
				}
			}

			// Расчет EV прибыли
			evProfit := sumVal * roi / 100
			log.Printf("EV Profit: %.2f (ROI: %.2f%%)", evProfit, roi)

			// Расчет прибыли
			var realProfit float64
			if result.Status == "win" {
				realProfit = sumVal * (coefVal - 1)
			} else if result.Status == "lose" {
				realProfit = -sumVal
			} else if result.Status == "return" || result.Status == "void" {
				realProfit = 0
			} else {
				log.Printf("Warning: bet status %s", result.Status)
				continue
			}

			// Обновляем информацию о прибыли в тестовой таблице
			err = s.Repo.UpdateBetProfits(bet.KeyOutcome, evProfit, realProfit)
			if err != nil {
				errorMsg := fmt.Sprintf("Failed to update profits in DB: %v", err)
				log.Print(errorMsg)
				if s.LogToTelegram && s.TelegramService != nil {
					s.TelegramService.SendMessage(errorMsg)
				}
			} else {
				debugMsg := fmt.Sprintf("Successfully updated profits for bet %s. EV: %.2f, Real: %.2f", bet.KeyOutcome, evProfit, realProfit)
				log.Print(debugMsg)
				if s.LogToTelegram && s.TelegramService != nil {
					s.TelegramService.SendMessage(debugMsg)
				}
			}
		}
		processedBets++
	}

	logMsg = fmt.Sprintf("Processed %d/%d bets", processedBets, totalBets)
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	return nil
}

func (s *BetService) ProcessTestRecentBets() error {
	logMsg := fmt.Sprintf("Processing test bets for the last day")
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	// Get test bets for the specified period
	bets, err := s.Repo.GetYesterdayTestBets()
	if err != nil {
		errMsg := fmt.Sprintf("Error getting recent test bets: %v", err)
		log.Print(errMsg)
		return err
	}

	logMsg = fmt.Sprintf("Found %d test bets for the last day", len(bets))
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	// Process each test bet
	totalBets := len(bets)
	processedBets := 0
	for _, bet := range bets {
		result, err := s.ProcessBet(bet)
		if err != nil {
			log.Printf("Error processing test bet %s: %v", bet.KeyOutcome, err)
		} else {
			log.Printf("Processed test bet %s: %s", bet.KeyOutcome, result.Status)

			// Обновляем информацию о прибыли для тестовых ставок
			if result.Status == "win" || result.Status == "lose" {
				// Извлекаем сумму ставки и коэффициент
				sumVal, ok := bet.Data["sum"].(float64)
				if !ok {
					log.Printf("Warning: couldn't extract sum from test bet %s", bet.KeyOutcome)
					sumVal = 0
				}
				coefVal, ok := bet.Data["coef"].(float64)
				if !ok {
					log.Printf("Warning: couldn't extract coefficient from test bet %s", bet.KeyOutcome)
					coefVal = 0
				}

				// Расчет прибыли
				var realProfit float64
				if result.Status == "win" {
					realProfit = sumVal * (coefVal - 1)
				} else if result.Status == "lose" {
					realProfit = -sumVal
				}

				// Обновляем информацию о прибыли в тестовой таблице
				err = s.Repo.UpdateTestBetProfits(bet.KeyOutcome, realProfit, realProfit)
				if err != nil {
					log.Printf("Error updating test bet profits %s: %v", bet.KeyOutcome, err)
				}

				processedBets++
			}
		}
	}

	logMsg = fmt.Sprintf("Processed %d/%d test bets", processedBets, totalBets)
	log.Print(logMsg)
	if s.LogToTelegram && s.TelegramService != nil {
		s.TelegramService.SendMessage(logMsg)
	}

	return nil
}
