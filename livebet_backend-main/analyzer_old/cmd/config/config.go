package config

import (
	"fmt"
	"livebets/analazer/pkg/utils"
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
	PostgresConfig `mapstructure:"postgres"`
	RedisConfig    `mapstructure:"redis"`
	PairsMatching  `mapstructure:"pairs_matching"`
	PriceStorage   `mapstructure:"price_storage"`
	MarketStorage  `mapstructure:"market_storage"`
	Port           `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DB       int64  `mapstructure:"db"`
	Password string `mapstructure:"password"`
	HashTTL  int64  `mapstructure:"hash_ttl"`
}

type PairsMatching struct {
	SendInterval             int `mapstructure:"send_interval"`
	UpdateKeysCacheInterval  int `mapstructure:"update_keys_cache_interval"`
	UpdatePairsCacheInterval int `mapstructure:"update_pairs_cache_interval"`
	ClearCacheInterval       int `mapstructure:"clear_cache_interval"`
	MatchDataTimeout         int `mapstructure:"match_data_timeout"`
	ReceiveWorkersCount      int `mapstructure:"receive_workers_count"`
	PairTimeout              int `mapstructure:"pair_timeout"`
}

type Port struct {
	Other    int `mapstructure:"other"`
	Pinnacle int `mapstructure:"pinnacle"`
	Server   int `mapstructure:"server"`
	Sender   int `mapstructure:"sender"`
}

type PriceStorage struct {
	ClearInterval int `mapstructure:"clear_interval"`
	DataTimeout   int `mapstructure:"data_timeout"`
}

type MarketStorage struct {
	ClearInterval int `mapstructure:"clear_interval"`
	DataTimeout   int `mapstructure:"data_timeout"`
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
