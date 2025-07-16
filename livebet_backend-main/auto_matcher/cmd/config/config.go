package config

import (
	"fmt"
	"livebets/auto_matcher/pkg/utils"
	"reflect"
	"strings"
	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var (
	once         sync.Once
	cachedConfig AppConfig
)

type AppConfig struct {
	PostgresConfig      `mapstructure:"postgres"`
	AnalyzerAPI         `mapstructure:"analyzer_api"`
	AnalyzerPrematchAPI AnalyzerAPI `mapstructure:"analyzer_pre_api"`
	AutoMatcherConfig   `mapstructure:"auto_matcher"`
	AIMatcherConfig     `mapstructure:"ai_matcher"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`
}

type AnalyzerAPI struct {
	URL          string `mapstructure:"url"`
	Timeout      int    `mapstructure:"timeout"`
	MatchDataURL string `mapstructure:"match_data_url"`
}

type AutoMatcherConfig struct {
	IntervalMatchingLeagues int `mapstructure:"interval_matching_leagues"`
	IntervalMatchingTeams   int `mapstructure:"interval_matching_teams"`
}

type AIMatcherConfig struct {
	ApiKey           string       `mapstructure:"api_key"`
	LiveInterval     int          `mapstructure:"live_interval"`
	PrematchInterval int          `mapstructure:"prematch_interval"`
	Claude           ClaudeConfig `mapstructure:"claude"`
}

type ClaudeConfig struct {
	Model     string `mapstructure:"model"`
	MaxTokens int    `mapstructure:"max_tokens"`
	Messages  struct {
		SystemLeaguesMsg string `mapstructure:"system_leagues_msg"`
		SystemTeamsMsg   string `mapstructure:"system_teams_msg"`
		BetcenterMsg     string `mapstructure:"betcenter_msg"`
		FonbetMsg        string `mapstructure:"fonbet_msg"`
		LadbrokesMsg     string `mapstructure:"ladbrokes_msg"`
		LobbetMsg        string `mapstructure:"lobbet_msg"`
		MaxbetMsg        string `mapstructure:"maxbet_msg"`
		SansabetMsg      string `mapstructure:"sansabet_msg"`
		SbbetMsg         string `mapstructure:"sbbet_msg"`
		StarCasinoMsg    string `mapstructure:"starcasino_msg"`
		UnibetMsg        string `mapstructure:"unibet_msg"`
	}
}

func (cfg PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

func ProvideAppMPConfig() (AppConfig, error) {
	var err error
	once.Do(func() {
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		viper.AddConfigPath("configs")
		viper.SetConfigName("common")
		viper.SetConfigType("yml")
		err = viper.ReadInConfig()
		if err != nil {
			return
		}

		BindEnvs(cachedConfig)

		hooks := viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(utils.DefaultDecodeHooks()...))
		err = viper.Unmarshal(&cachedConfig, hooks)
		if err != nil {
			return
		}
	})

	return cachedConfig, err
}

func BindEnvs(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
