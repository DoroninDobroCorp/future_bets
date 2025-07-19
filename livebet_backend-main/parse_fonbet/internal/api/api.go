package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"livebets/parse_fonbet/cmd/config"
	"livebets/parse_fonbet/internal/entity"
	"log"
	"net/http"
	"net/url"
	"time"
)

type FonbetAPI struct {
	cfg    config.FonbetConfig
	client *http.Client
}

func NewFonbetAPI(cfg config.FonbetConfig) *FonbetAPI {
	return &FonbetAPI{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (api *FonbetAPI) GetAllData() (*entity.Requested, error) {
	start := time.Now()

	req, err := http.NewRequest(
		http.MethodGet,
		api.cfg.Url+api.cfg.MatchesUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	proxy, _ := url.Parse(api.cfg.ProxyUrl)
	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	api.client.Transport = transport

	query := req.URL.Query()
	query.Add("lang", "en")
	query.Add("scopeMarket", "1600")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", api.cfg.Url)
	req.Header.Set("Origin", "https://www.fon.bet")
	req.Header.Set("Referer", "https://www.fon.bet/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

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

	var requested entity.Requested
	if err := json.NewDecoder(&bodyBuffer).Decode(&requested); err != nil {
		return nil, err
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Время получения данных о всех матчах: %s", elapsed)

	return &requested, nil
}

func (api *FonbetAPI) GetMatchODDS(matchId int64) (*entity.Requested, error) {
	start := time.Now()

	req, err := http.NewRequest(
		http.MethodGet,
		api.cfg.Url+api.cfg.ODDSUrl,
		nil,
	)
	if err != nil {
		return nil, err
	}

	proxy, _ := url.Parse(api.cfg.ProxyUrl)
	transport := &http.Transport{Proxy: http.ProxyURL(proxy)}
	api.client.Transport = transport

	query := req.URL.Query()
	query.Add("lang", "en")
	query.Add("version", "0")
	query.Add("eventId", fmt.Sprintf("%d", matchId))
	query.Add("scopeMarket", "1600")
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.8,en-US;q=0.5,en;q=0.3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", api.cfg.Url)
	req.Header.Set("Origin", "https://www.fon.bet")
	req.Header.Set("Referer", "https://www.fon.bet/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")

	resp, err := api.client.Do(req)
	if err != nil {
		log.Printf("[ERROR0] %s", err.Error())
		return nil, err
	}

	body := resp.Body

	if resp.Header.Get("Content-Encoding") == "gzip" {
		body, err = gzip.NewReader(resp.Body)
		if err != nil {
			log.Printf("[ERROR1] %s (%v)", err.Error(), resp.Header)
			return nil, err
		}
		defer body.Close()
	}

	var bodyBuffer bytes.Buffer
	_, err = bodyBuffer.ReadFrom(body)
	if err != nil {
		log.Printf("[ERROR2] %s", err.Error())
		return nil, err
	}

	var requested entity.Requested
	if err := json.NewDecoder(&bodyBuffer).Decode(&requested); err != nil {
		log.Printf("[ERROR3] %s", err.Error())
		return nil, err
	}

	elapsed := time.Since(start)
	log.Printf("[INFO] Время получения данных о всех кэфах: %s", elapsed)

	return &requested, nil
}
