package api

import (
	"encoding/json"
	"livebets/auto_matcher/cmd/config"
	"livebets/auto_matcher/internal/entity"
	"net/http"
	"time"
)

type AnalizerPrematchAPI struct {
	cfg    config.AnalyzerAPI
	client *http.Client
}

func NewAnalizerPrematchAPI(cfg config.AnalyzerAPI) *AnalizerPrematchAPI {
	transport := &http.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(cfg.Timeout),
	}

	return &AnalizerPrematchAPI{
		cfg:    cfg,
		client: client,
	}
}

func (a *AnalizerPrematchAPI) GetOnlineMatchData() ([]entity.MatchData, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		a.cfg.URL+a.cfg.MatchDataURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []entity.MatchData
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
