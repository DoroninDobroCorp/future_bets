package api

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/andybalholm/brotli"
	"golang.org/x/net/http2"
	"io"
	"livebets/parse_starcasino/cmd/config"
	"livebets/parse_starcasino/internal/entity"
	"livebets/shared"
	"log"
	"net/http"
	"strconv"
	"time"
)

type API struct {
	cfg    config.APIConfig
	client *http.Client
}

func New(cfg config.APIConfig) *API {
	transport := &http2.Transport{}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Second * time.Duration(cfg.Timeout),
	}

	return &API{
		cfg:    cfg,
		client: client,
	}
}

func (api *API) GetLiveEvents(sportName shared.SportName) ([]*entity.Event, error) {
	start := time.Now()

	sportId := ""
	switch sportName {
	case shared.SOCCER:
		sportId = "66" // Football
	case shared.TENNIS:
		sportId = "68"
	default:
		return nil, fmt.Errorf("incorrect events sport name: %s", sportName)
	}

	req, err := http.NewRequest(http.MethodGet, api.cfg.Url+api.cfg.Live.EventsUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("culture", "en-GB")
	query.Add("timezoneOffset", "-60")
	query.Add("integration", "starcasino.be")
	query.Add("deviceType", "1")
	query.Add("numFormat", "en-GB")
	query.Add("countryCode", "BE")
	query.Add("sportId", sportId)
	req.URL.RawQuery = query.Encode()

	req.Header = http.Header{
		"Accept":             {"*/*"},
		"Accept-Encoding":    {"gzip, deflate, br, zstd"},
		"Accept-Language":    {"en-GB;q=0.9"},
		"Origin":             {"https://starcasinosport.be"},
		"Priority":           {"u=1, i"},
		"Referer":            {"https://starcasinosport.be/"},
		"Sec-Ch-Ua":          {`"Not A(Brand";v="8", "Chromium";v="132", "Microsoft Edge";v="132"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"cross-site"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"},
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var eventsData entity.EventsData
	if err := json.Unmarshal(*body, &eventsData); err != nil {
		return nil, err
	}

	log.Printf("[INFO] Время получения %s данных %s", sportName, time.Since(start))

	return eventsData.Events, nil
}

func (api *API) GetLiveOneEvent(eventId int64) (*entity.Match, error) {
	req, err := http.NewRequest(http.MethodGet, api.cfg.Url+api.cfg.Live.OddsUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("culture", "en-GB")
	query.Add("timezoneOffset", "-60")
	query.Add("integration", "starcasino.be")
	query.Add("deviceType", "1")
	query.Add("numFormat", "en-GB")
	query.Add("countryCode", "BE")
	query.Add("eventId", strconv.FormatInt(eventId, 10))
	req.URL.RawQuery = query.Encode()

	req.Header = http.Header{
		"Accept":             {"*/*"},
		"Accept-Encoding":    {"gzip, deflate, br, zstd"},
		"Accept-Language":    {"en-GB;q=0.9"},
		"Origin":             {"https://starcasinosport.be"},
		"Priority":           {"u=1, i"},
		"Referer":            {"https://starcasinosport.be/"},
		"Sec-Ch-Ua":          {`"Not A(Brand";v="8", "Chromium";v="132", "Microsoft Edge";v="132"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"cross-site"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"},
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	match := &entity.Match{}
	if err := json.Unmarshal(*body, match); err != nil {
		return nil, err
	}

	return match, nil
}

func getResponseBody(response *http.Response) (*[]byte, error) {
	var reader io.Reader

	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		var err error
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	case "br":
		reader = brotli.NewReader(response.Body)
	default:
		reader = response.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func (api *API) GetPreMatchPage(sportName shared.SportName, pageNumber int, eventsPerPage int) (*entity.PreMatchData, error) {
	start := time.Now()

	sportId := ""
	switch sportName {
	case shared.SOCCER:
		sportId = "66" // Football
	case shared.TENNIS:
		sportId = "68"
	default:
		return nil, fmt.Errorf("incorrect events sport name: %s", sportName)
	}

	req, err := http.NewRequest(http.MethodGet, api.cfg.Url+api.cfg.Prematch.EventsUrl, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("culture", "en-GB")
	query.Add("timezoneOffset", "-60")
	query.Add("integration", "starcasino.be")
	query.Add("deviceType", "1")
	query.Add("numFormat", "en-GB")
	query.Add("countryCode", "BE")
	query.Add("sportId", sportId)
	query.Add("eventCount", strconv.Itoa(eventsPerPage))
	if pageNumber > 1 {
		query.Add("page", strconv.Itoa(pageNumber))
	}
	req.URL.RawQuery = query.Encode()

	req.Header = http.Header{
		"Accept":             {"*/*"},
		"Accept-Encoding":    {"gzip, deflate, br, zstd"},
		"Accept-Language":    {"en-GB;q=0.9"},
		"Origin":             {"https://starcasinosport.be"},
		"Priority":           {"u=1, i"},
		"Referer":            {"https://starcasinosport.be/"},
		"Sec-Ch-Ua":          {`"Not A(Brand";v="8", "Chromium";v="132", "Microsoft Edge";v="132"`},
		"Sec-Ch-Ua-Mobile":   {"?0"},
		"Sec-Ch-Ua-Platform": {`"Windows"`},
		"Sec-Fetch-Dest":     {"empty"},
		"Sec-Fetch-Mode":     {"cors"},
		"Sec-Fetch-Site":     {"cross-site"},
		"User-Agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0"},
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	body, err := getResponseBody(resp)
	if err != nil {
		return nil, err
	}

	preMatchData := &entity.PreMatchData{}
	if err := json.Unmarshal(*body, preMatchData); err != nil {
		return nil, err
	}

	log.Printf("[INFO] Время получения %s данных %s. Страница %d Матчей %d", sportName, time.Since(start), pageNumber, len(preMatchData.Events))

	return preMatchData, nil
}
