package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"livebets/parse_lobbet/cmd/config"
	"livebets/parse_lobbet/internal/entity"
	"livebets/shared"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type API struct {
	cfg    config.APIConfig
	client *http.Client
}

func New(cfg config.APIConfig) *API {
	transport := &http.Transport{}
	if cfg.Proxy != "" {
		proxyURL, err := url.Parse(cfg.Proxy)
		if err != nil {
			log.Fatalf("Bad proxy URL: %s", cfg.Proxy)
		} else {
			transport.Proxy = http.ProxyURL(proxyURL)
			log.Printf("Transport proxy: %s", proxyURL.String())
		}
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(cfg.Timeout),
	}

	return &API{
		cfg:    cfg,
		client: client,
	}
}

const preMatchTimeDiff = 48 * time.Hour

func (api *API) GetAllMatches(filterStatus int64) ([]*entity.Match, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		api.cfg.Url+api.cfg.Live.EventsUrl,
		bytes.NewBuffer([]byte("{}")),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	if resp.Body == nil {
		return nil, err // err ?
	}

	encodedBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	defer encodedBody.Close()

	body, err := io.ReadAll(encodedBody)
	if err != nil {
		return nil, err
	}

	var apiResponse entity.ResponseMatchData
	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	timeNow := time.Now()

	sportFilter := makeSportFilter(api.cfg.SportConfig)

	matches := make([]*entity.Match, 0, 128)

	for _, match := range apiResponse.Live.Matches {
		if match.LiveStatus != filterStatus {
			continue
		}

		if _, ok := sportFilter[match.SportLetter]; !ok {
			continue
		}

		if filterStatus == 1 { // Live
			if len(match.Bets) == 0 {
				continue
			}
		} else { // PreMatch
			// Check match start time
			startTime := time.Unix(match.TimeStamp/1000, 0)

			// Matches that start no later than 48 hours
			diff := -timeNow.Sub(startTime) // -(minus) before timeNow.Sub() is to diff will be > 0
			if diff < 0 || diff > preMatchTimeDiff {
				continue
			}

		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (api *API) GetMatchOdds(match *entity.Match) error {
	odds_url := strings.ReplaceAll(api.cfg.Prematch.OddsUrl, "{matchId}", fmt.Sprintf("%d", match.ID))
	req, err := http.NewRequest(
		http.MethodGet,
		api.cfg.Url+odds_url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	if resp.Body == nil {
		return nil
	}

	encodedBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer encodedBody.Close()

	body, err := io.ReadAll(encodedBody)
	if err != nil {
		return err
	}

	var preMatchData entity.ResponsePrematchData
	if err = json.Unmarshal(body, &preMatchData); err != nil {
		return err
	}

	var bets []entity.Bet

	for _, prematchBet := range preMatchData.Bets {
		var picks []entity.Pick
		for _, tip := range prematchBet.Tips {
			if tip.Value != .0 {
				picks = append(picks, entity.Pick{
					Caption:      tip.Name,
					OddValue:     tip.Value,
					SpecialValue: prematchBet.HandicapParamValue,
				})
			}
		}

		if len(picks) > 0 {
			bets = append(bets, entity.Bet{
				LiveBetCaption: prematchBet.Description,
				Picks:          picks,
			})
		}
	}

	match.Bets = bets

	return nil
}

func makeSportFilter(sportConfig config.SportConfig) map[string]shared.SportName {
	filter := make(map[string]shared.SportName)

	if sportConfig.Football {
		filter["S"] = shared.SOCCER
	}

	if sportConfig.Tennis {
		filter["T"] = shared.TENNIS
	}

	if sportConfig.Basketball {
		filter["B"] = shared.BASKETBALL
	}

	return filter
}
