package api

import (
	"net/http"
	"time"
)

type ParserAPI struct {
	client *http.Client
}

func NewParserAPI() *ParserAPI {
	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(5),
	}

	return &ParserAPI{
		client: client,
	}
}

func (a *ParserAPI) GetOnlineMatchData(url string) (int, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
