package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"livebets/parse_maxbet/cmd/config"
	"livebets/parse_maxbet/internal/entity"
	"log"
	"net/http"
	"strings"
	"time"
)

type API struct {
	cfg    config.APIConfig
	client *http.Client
}

func New(cfg config.APIConfig) *API {
	client := &http.Client{
		Timeout: time.Second * time.Duration(cfg.Timeout),
	}

	return &API{
		cfg:    cfg,
		client: client,
	}
}

func (api *API) GetLiveData() ([]entity.RequestedData, error) {
	start := time.Now()
	url := "https://www.maxbet.me/live/events/sr_ME"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error api.client.Do(): %v", err)
	}

	body, err := getResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("error getResponseBody(): %v", err)
	}

	bodyStr := string(*body)
	splits := strings.Split(bodyStr, "data:")
	batchesCount := len(splits)

	results := make([]entity.RequestedData, 0, batchesCount-2)

	for i, splitted := range splits {
		if i == 0 || i == batchesCount-1 {
			continue
		}

		result := entity.RequestedData{}

		if err := json.Unmarshal([]byte(splitted), &result); err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Время получения данных: %s", elapsed)

	return results, nil
}

func getResponseBody(response *http.Response) (*[]byte, error) {
	var reader io.Reader
	var err error

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	//case "br":
	//	reader = brotli.NewReader(response.Body)
	//case "zstd":
	//	reader, err = zstd.NewReader(response.Body)
	//	if err != nil {
	//		return nil, err
	//	}
	default:
		reader = response.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
