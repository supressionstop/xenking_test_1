package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App        App        `mapstructure:"app"`
		Log        Log        `mapstructure:"log"`
		HTTPServer HTTPServer `mapstructure:"http_server"`
		DB         DB         `mapstructure:"db"`
		Provider   Provider   `mapstructure:"provider"`
		Workers    []Worker   `mapstructure:"workers"`
		GRPCServer GrpcServer `mapstructure:"grpc_server"`
	}

	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	HTTPServer struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}

	DB struct {
		URL string `mapstructure:"url"`
	}

	Provider struct {
		BaseURL     string        `example:"http://localhost:8080" mapstructure:"base_url"`
		HTTPTimeout time.Duration `example:"5s"                    mapstructure:"http_timeout"`
	}

	Worker struct {
		Sport        string        `mapstructure:"sport"`
		PollInterval time.Duration `mapstructure:"poll_interval"`
	}

	GrpcServer struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
)

func MustSetup(environment string) (*Config, error) {
	// 1st priority - env
	// DB_URL (env) ->  Config.DB.URL (struct)
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))
	viper.AutomaticEnv()
	// 2nd priority - json
	if environment == "" {
		viper.SetConfigName("config")
	} else {
		viper.SetConfigName(fmt.Sprintf("config_%s", strings.ToLower(environment)))
	}

	viper.SetConfigType("json")

	configFolderPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configFolderPath = filepath.Join(configFolderPath, "config")
	viper.AddConfigPath(configFolderPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{} //nolint:exhaustruct

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
