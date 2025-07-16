package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"livebets/results/internal/api/model"
	"livebets/results/internal/entity"
)

var SPORTS = map[string]int{
	"Soccer":     29,
	"Tennis":     33,
	"Basketball": 4,
}

type PinnacleResponse struct {
	Leagues []PinnacleLeague `json:"leagues"`
}

type PinnacleLeague struct {
	Events []model.Event `json:"events"`
}

type PinnacleService struct {
	Login          string
	Password       string
	ProxyEnabled   bool
	ProxyURL       string
	TelegramLogger *TelegramService
	CachedResponse map[int]PinnacleResponse
}

func NewPinnacleService(login, password string, proxyEnabled bool, proxyURL string) *PinnacleService {
	return &PinnacleService{
		Login:          login,
		Password:       password,
		ProxyEnabled:   proxyEnabled,
		ProxyURL:       proxyURL,
		CachedResponse: make(map[int]PinnacleResponse),
	}
}

func (s *PinnacleService) SetTelegramLogger(telegramService *TelegramService) {
	s.TelegramLogger = telegramService
}

func (s *PinnacleService) CallFixtureSettled(bet *entity.LogBetAccept) (*model.Fixture, error) {
	// Check cache first
	fixtureKey := bet.KeyMatch
	if strings.Contains(fixtureKey, ":") {
		fixtureKey = strings.Split(fixtureKey, ":")[0]
	}

	// Collect data for API
	var matchId int
	if pair, ok := bet.Data["pair"].(map[string]interface{}); ok {
		if first, ok := pair["first"].(map[string]interface{}); ok {
			if id, ok := first["matchId"].(string); ok {
				matchId, _ = strconv.Atoi(id)
			}
		}
	}

	var sportId int
	if pair, ok := bet.Data["pair"].(map[string]interface{}); ok {
		if sportName, ok := pair["sportName"].(string); ok {
			sportId = SPORTS[sportName]
		}
	}

	// Log the API call
	logMsg := fmt.Sprintf("Calling Pinnacle API for fixture: %s", fixtureKey)
	log.Print(logMsg)

	// Call the API endpoint
	logMsg = fmt.Sprintf("API parameters: matchId=%d, sportId=%d", matchId, sportId)
	log.Print(logMsg)

	fixture, err := s.callFixturesEndpoint(matchId, sportId)
	if err != nil {
		return nil, err
	}

	return fixture, nil
}

func (s *PinnacleService) callFixturesEndpoint(matchId, sportId int) (*model.Fixture, error) {
	// get data from cache if it exists
	if _, ok := s.CachedResponse[sportId]; ok {
		matchedEvent := findEvent(s.CachedResponse[sportId], matchId)
		if matchedEvent == nil {
			errMsg := fmt.Sprintf("No fixture found matching matchId=%d", matchId)
			log.Print(errMsg)
			return nil, fmt.Errorf("fixture not found: %s", errMsg)
		}

		fixture := &model.Fixture{
			Periods: matchedEvent.Periods,
		}

		logPeriods := ""
		for _, period := range fixture.Periods {
			logPeriods += fmt.Sprintf("Period %d: %d-%d, ",
				period.Number, period.Team1Score, period.Team2Score)
		}
		logMsg := fmt.Sprintf("Found matching fixture: matchId=%d, Scores: %s", matchId, logPeriods)
		log.Print(logMsg)

		return fixture, nil
	}

	// Prepare URL with parameters
	settledUrl := fmt.Sprintf("http://65.109.115.96:5556/fixtures/settled?sportId=%d", sportId)

	// Create HTTP client with optional proxy
	var client *http.Client
	if s.ProxyEnabled && s.ProxyURL != "" {
		proxyURL := s.ProxyURL
		client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(mustParseURL(proxyURL)),
			},
			Timeout: 30 * time.Second,
		}
		log.Printf("Using proxy: %s", proxyURL)
	} else {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
		log.Print("No proxy used")
	}

	// Create request
	req, err := http.NewRequest("GET", settledUrl, nil)
	if err != nil {
		errMsg := fmt.Sprintf("Error creating request: %v", err)
		log.Print(errMsg)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("token", "VCBSIHBpbm5hY2xldHJhbnNsYXRvcgo=")

	// Make the request
	logMsg := fmt.Sprintf("Sending request to URL: %s", settledUrl)
	log.Print(logMsg)

	resp, err := client.Do(req)
	if err != nil {
		errMsg := fmt.Sprintf("Error making request: %v", err)
		log.Print(errMsg)
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Sprintf("API returned non-200 status: %d", resp.StatusCode)
		log.Print(errMsg)
		return nil, fmt.Errorf("API error: %s", errMsg)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Error reading response body: %v", err)
		log.Print(errMsg)
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	// Parse JSON response
	var pinnacleResp PinnacleResponse
	err = json.Unmarshal(body, &pinnacleResp)
	if err != nil {
		errMsg := fmt.Sprintf("Error parsing JSON response: %v", err)
		log.Print(errMsg)
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	s.CachedResponse[sportId] = pinnacleResp

	// Log response structure
	logMsg = fmt.Sprintf("Received response with %d leagues", len(pinnacleResp.Leagues))
	log.Print(logMsg)

	// Find the fixture that matches matchId
	matchedEvent := findEvent(pinnacleResp, matchId)
	if matchedEvent == nil {
		errMsg := fmt.Sprintf("No fixture found matching matchId=%d", matchId)
		log.Print(errMsg)
		return nil, fmt.Errorf("fixture not found: %s", errMsg)
	}

	// Create and return the fixture
	fixture := &model.Fixture{
		Periods: matchedEvent.Periods,
	}

	// Log the matched event score
	logPeriods := ""
	for _, period := range fixture.Periods {
		logPeriods += fmt.Sprintf("Period %d: %d-%d, ",
			period.Number, period.Team1Score, period.Team2Score)
	}
	logMsg = fmt.Sprintf("Found matching fixture: matchId=%d, Scores: %s", matchId, logPeriods)
	log.Print(logMsg)

	return fixture, nil
}

func findEvent(pinnacleResp PinnacleResponse, matchId int) *model.Event {
	for _, league := range pinnacleResp.Leagues {
		for _, event := range league.Events {
			if event.ID == matchId {
				return &event
			}
		}
	}
	return nil
}

func (s *PinnacleService) ClearCached() {
	for k, _ := range s.CachedResponse {
		delete(s.CachedResponse, k)
	}
}

func mustParseURL(rawURL string) *url.URL {
	proxyURL, err := url.Parse(rawURL)
	if err != nil || proxyURL == nil {
		log.Print("Error parsing proxy URL, using direct connection")
		return nil
	}
	return proxyURL
}
