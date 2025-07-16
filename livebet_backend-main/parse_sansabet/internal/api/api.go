package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"livebets/parse_sansabet/cmd/config"
	"livebets/parse_sansabet/internal/entity"
	"log"
	"net/http"
	"time"
)

var (
	slidAll int64 = 0
)

type SansabetAPI struct {
	cfg    config.SansabetConfig
	client *http.Client
}

func NewSansabetAPI(cfg config.SansabetConfig) *SansabetAPI {
	return &SansabetAPI{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (api *SansabetAPI) GetAllMatches() (*[]entity.RequestedEvent, error) {
	start := time.Now()

	req, err := http.NewRequest(
		http.MethodGet,
		api.cfg.Url+api.cfg.MatchesUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("SLID", fmt.Sprintf("%d", slidAll))
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Origin", "https://sansabet.com")
	req.Header.Set("Referer", "https://sansabet.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	encodedBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer encodedBody.Close()

	var bodyBuffer bytes.Buffer
	_, err = bodyBuffer.ReadFrom(encodedBody)
	if err != nil {
		return nil, err
	}

	var events []entity.RequestedEvent
	if err := json.NewDecoder(&bodyBuffer).Decode(&events); err != nil {
		return nil, err
	}

	for _, event := range events {
		if event.H.SLID > slidAll {
			slidAll = event.H.SLID
		}
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Время получения данных о всех матчах: %s", elapsed)

	return &events, nil
}

func (api *SansabetAPI) parseOdds(matchId int64) (*entity.EventOdds, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		api.cfg.Url+api.cfg.ODDSUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("SLID", fmt.Sprintf("%d", 0))
	query.Add("ParIDs", fmt.Sprintf("%d", matchId))
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Origin", "https://sansabet.com")
	req.Header.Set("Referer", "https://sansabet.com/")
	req.Header.Set("Sec-Ch-Ua", "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Linux\"")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64)")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	encodedBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer encodedBody.Close()

	var bodyBuffer bytes.Buffer
	_, err = bodyBuffer.ReadFrom(encodedBody)
	if err != nil {
		return nil, err
	}

	var result []entity.EventOdds
	if err := json.NewDecoder(&bodyBuffer).Decode(&result); err != nil {
		return nil, err
	}

	return &result[0], nil
}

func (api *SansabetAPI) GetAllMatchesODDS(matchIds []int64) (*[]entity.EventOdds, error) {
	start := time.Now()

	var result []entity.EventOdds

	for _, matchId := range matchIds {
		parsedOdds, _ := api.parseOdds(matchId)
		result = append(result, *parsedOdds)
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Время получения данных о всех кэфах: %s", elapsed)

	return &result, nil
}
