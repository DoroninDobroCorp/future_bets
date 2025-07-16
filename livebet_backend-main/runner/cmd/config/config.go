package config

import (
	"fmt"
	"livebets/runner/pkg/utils"
	"os"
	"reflect"
	"strconv"
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
	Bookmakers    []BookmakerConfig `mapstructure:"bookmaker"`
	CommandConfig `mapstructure:"command"`
	StatusConfig  `mapstructure:"status"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db"`
	SSLMode  string `mapstructure:"sslmode"`
}

func (cfg PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}

type BookmakerConfig struct {
	Replicas     int    `mapstructure:"replicas"`
	ReplicasName string `mapstructure:"replicas_name"`
	Name         string `mapstructure:"name"`
	Path         string `mapstructure:"path"`
	API          string `mapstructure:"api"`
}

type CommandConfig struct {
	Port    string `mapstructure:"port"`
	EnvPath string `mapstructure:"env_path"`
}

type StatusConfig struct {
	Interval int `mapstructure:"interval"`
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

		if bookmakers := bindEnvBookmakers(); bookmakers != nil {
			cachedConfig.Bookmakers = bookmakers
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

func bindEnvBookmakers() []BookmakerConfig {
	var rt []BookmakerConfig
	for i := 0; i < 200; i++ {
		envName, ok := os.LookupEnv(fmt.Sprintf("BOOKMAKER_%d_NAME", i))
		if !ok {
			continue
		}
		envPath, ok := os.LookupEnv(fmt.Sprintf("BOOKMAKER_%d_PATH", i))
		if !ok {
			continue
		}
		envAPI, ok := os.LookupEnv(fmt.Sprintf("BOOKMAKER_%d_API", i))
		if !ok {
			continue
		}
		envPeplicasName, ok := os.LookupEnv(fmt.Sprintf("BOOKMAKER_%d_REPLICAS_NAME", i))
		if !ok {
			continue
		}
		envReplicas, ok := os.LookupEnv(fmt.Sprintf("BOOKMAKER_%d_REPLICAS", i))
		if !ok {
			continue
		}
		replicas, err := strconv.Atoi(envReplicas)
		if err != nil {
			continue
		}
		rt = append(rt, BookmakerConfig{
			Replicas:     replicas,
			ReplicasName: envPeplicasName,
			Name:         envName,
			Path:         envPath,
			API:          envAPI,
		})
	}

	return rt
}
