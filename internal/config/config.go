package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

	cfg := &Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	if err := cfg.setupWorkersFromEnvs(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) MaxWorkerInterval() time.Duration {
	var interval time.Duration
	for _, w := range cfg.Workers {
		if w.PollInterval > interval {
			interval = w.PollInterval
		}
	}

	return interval
}

func (cfg *Config) setupWorkersFromEnvs() error {
	const prefix = "WORKERS_"
	var workerEnvs []string
	envs := os.Environ()
	for _, env := range envs {
		if strings.HasPrefix(env, prefix) {
			workerEnvs = append(workerEnvs, env)
		}
	}

	if len(workerEnvs) == 0 {
		return nil
	}

	workersCfg := make([]Worker, len(workerEnvs)/2)
	for _, env := range workerEnvs {
		pair := strings.SplitN(env, "=", 2)
		name := pair[0]
		value := pair[1]

		nameParts := strings.Split(name, "_")
		if len(nameParts) < 2 {
			return fmt.Errorf("cant parse worker idx: %s", name)
		}
		idx, err := strconv.Atoi(nameParts[1])
		if err != nil {
			return fmt.Errorf("cant parse worker idx: %s", name)
		}

		if len(nameParts) == 3 && nameParts[2] == "SPORT" {
			workersCfg[idx].Sport = value
		}
		if len(nameParts) == 4 && nameParts[2] == "POLL" && nameParts[3] == "INTERVAL" {
			pollInterval, err := time.ParseDuration(value)
			if err != nil {
				return fmt.Errorf("failed to parse worker poll interval: %q", value)
			}
			workersCfg[idx].PollInterval = pollInterval
		}
	}

	cfg.Workers = workersCfg

	return nil
}
