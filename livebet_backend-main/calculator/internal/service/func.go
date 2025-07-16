package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"livebets/calculator/internal/entity"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func removeSpecialChars(input string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_<>-]`)
	return re.ReplaceAllString(input, "")
}

func replaceDotsInFileName(fileName string) string {
	return strings.ReplaceAll(fileName, ".", "-")
}

func normalizeTime(recordTime time.Time, minutes, seconds int, isLive bool) time.Time {
	resultTime := time.Date(
		recordTime.Year(),
		recordTime.Month(),
		recordTime.Day(),
		recordTime.Hour(),
		minutes,
		seconds,
		recordTime.Nanosecond(),
		recordTime.Location(),
	)

	if isLive {
		resultTime = resultTime.Add(5 * time.Second)
	} else {
		resultTime = resultTime.Add(18 * time.Minute)
	}

	return resultTime
}

func getPriceForSecond(priceRecords *entity.ResponsePriceRecords, recordTime time.Time, minutes, seconds int, isLive bool, coef float64) float64 {
	searchTime := normalizeTime(recordTime, minutes, seconds, isLive)

	for _, record := range priceRecords.Records {
		if record.Second.CreatedAt == searchTime {
			return record.Second.Score
		}
	}

	return coef
}

func sendMissedBet(pair entity.PairOneOutcome, keyMatch string) error {
	url := "http://188.253.24.91:7020/missed_bet"

	requestBody := entity.MissedBet{
		KeyMatch: keyMatch,
		Pair:     pair,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
