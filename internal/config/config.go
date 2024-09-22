package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

type (
	Config struct {
		App        App        `mapstructure:"app"`
		Log        Log        `mapstructure:"log"`
		HttpServer HttpServer `mapstructure:"http_server"`
	}

	App struct {
		Name        string `mapstructure:"name"`
		Environment string `mapstructure:"environment"`
	}

	Log struct {
		Level string `mapstructure:"level"`
	}

	HttpServer struct {
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

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	//appRootPath := filepath.Join(b, "../..")
	//setPathsFromRoot(appRootPath, cfg)

	return cfg, nil
}

func setPathsFromRoot(projectRoot string, config *Config) {
	// TODO: migrations
}
