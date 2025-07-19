package config

import (
	"livebets/parse_fonbet/utils"
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
	FonbetConfig `mapstructure:"fonbet"`
	SenderConfig `mapstructure:"sender"`
	Port         string `mapstructure:"port"`
}

type FonbetConfig struct {
	FonbetAPIConfig `mapstructure:"api"`
}

type FonbetAPIConfig struct {
	Url           string `mapstructure:"url"`
	Timeout       int    `mapstructure:"timeout"`
	MatchesUrl    string `mapstructure:"matches_url"`
	ODDSUrl       string `mapstructure:"odds_url"`
	ProxyUrl      string `mapstructure:"proxy_url"`
	IntervalMatch int    `mapstructure:"interval_match"`
	IntervalODDS  int    `mapstructure:"interval_odds"`
}

type SenderConfig struct {
	Url string `mapstructure:"url"`
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
