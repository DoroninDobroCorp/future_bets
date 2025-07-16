package api

import (
	"encoding/csv"
	"fmt"
	"livebets/results/internal/entity"
	"log"
	"math/rand"
	"os"
	"time"
)

type NewTableRow struct {
	KeyMatch   string
	KeyOutcome string
	MatchName  string
	OutcomeStr string
	Margin     float64
	ROI        float64
	CoefFirst  float64
	CoefSecond float64
	EVProfit   float64
	RealProfit float64
}

var config = struct {
	TableHeaders []string
	Sum          float64
}{
	TableHeaders: []string{
		"key_match", "key_outcome",
		"match_name", "outcome",
		"margin", "roi", "coef_first", "coef_second",
		"ev_profit", "real_profit",
	},
	Sum: 100,
}

var alreadyProcessed = make(map[string]bool)

func ProcessTable(rows []entity.LogBetAccept) []NewTableRow {
	// Перемешиваем строки
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(rows), func(i, j int) {
		rows[i], rows[j] = rows[j], rows[i]
	})

	var newTable []NewTableRow

	for _, row := range rows {
		// Проверяем, обрабатывали ли мы уже эту строку
		if alreadyProcessed[row.KeyMatch] {
			continue
		}

		// Извлекаем данные из JSON
		pair, ok := row.Data["pair"].(map[string]interface{})
		if !ok {
			log.Printf("Не удалось извлечь pair из данных для key_match: %s", row.KeyMatch)
			continue
		}

		outcome, ok := pair["outcome"].(map[string]interface{})
		if !ok {
			log.Printf("Не удалось извлечь outcome из данных для key_match: %s", row.KeyMatch)
			continue
		}

		score1, ok := outcome["score1"].(map[string]interface{})
		if !ok {
			log.Printf("Не удалось извлечь score1 из данных для key_match: %s", row.KeyMatch)
			continue
		}

		score2, ok := outcome["score2"].(map[string]interface{})
		if !ok {
			log.Printf("Не удалось извлечь score2 из данных для key_match: %s", row.KeyMatch)
			continue
		}

		first, ok := pair["first"].(map[string]interface{})
		if !ok {
			log.Printf("Не удалось извлечь first из данных для key_match: %s", row.KeyMatch)
			continue
		}

		// Преобразуем значения
		margin, _ := outcome["margin"].(float64)
		roi, _ := outcome["roi"].(float64)
		coef1, _ := score1["value"].(float64)
		coef2, _ := score2["value"].(float64)

		if roi > 12 {
			continue
		}

		coefVal := row.Data["coef"].(float64)
		outcomeStr := outcome["outcome"].(string)
		homeName := first["homeName"].(string)
		awayName := first["awayName"].(string)
		rowRealProfit := *row.RealProfit

		// Рассчитываем EVProfit
		evProfit := config.Sum * roi / 100

		// Рассчитываем RealProfit

		var realProfit float64
		if rowRealProfit > 0 {
			realProfit = config.Sum * (coefVal - 1)

		} else if rowRealProfit < 0 {
			realProfit = -config.Sum

		} else {
			realProfit = 0
		}

		// Создаем новую строку
		newRow := NewTableRow{
			KeyMatch:   row.KeyMatch,
			KeyOutcome: row.KeyOutcome,
			MatchName:  fmt.Sprintf("%s - %s", homeName, awayName),
			OutcomeStr: outcomeStr,
			Margin:     margin,
			ROI:        roi,
			CoefFirst:  coef1,
			CoefSecond: coef2,
			EVProfit:   evProfit,
			RealProfit: realProfit,
		}

		newTable = append(newTable, newRow)
		alreadyProcessed[row.KeyMatch] = true
	}

	return newTable
}

func SaveToCSV(rows []NewTableRow) error {
	file, err := os.Create("logs/statistic.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовки
	if err := writer.Write(config.TableHeaders); err != nil {
		return err
	}

	// Переменные для суммирования
	var totalCount int
	var totalEVProfit float64
	var totalRealProfit float64

	// Записываем данные и считаем суммы
	for _, row := range rows {
		record := []string{
			row.KeyMatch,
			row.KeyOutcome,
			row.MatchName,
			row.OutcomeStr,
			fmt.Sprintf("%.2f", row.Margin),
			fmt.Sprintf("%.2f", row.ROI),
			fmt.Sprintf("%.2f", row.CoefFirst),
			fmt.Sprintf("%.2f", row.CoefSecond),
			fmt.Sprintf("%.2f", row.EVProfit),
			fmt.Sprintf("%.2f", row.RealProfit),
		}
		if err := writer.Write(record); err != nil {
			return err
		}

		totalCount++
		totalEVProfit += row.EVProfit
		totalRealProfit += row.RealProfit
	}

	if err := writer.Write([]string{}); err != nil {
		return err
	}

	summaryHeaders := []string{
		fmt.Sprintf("All bets: %d", totalCount),
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"Summary evProfit:",
		"Summary realProfit:",
	}
	if err := writer.Write(summaryHeaders); err != nil {
		return err
	}

	summary := []string{
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		fmt.Sprintf("%.2f", totalEVProfit),
		fmt.Sprintf("%.2f", totalRealProfit),
	}
	if err := writer.Write(summary); err != nil {
		return err
	}

	return nil
}
