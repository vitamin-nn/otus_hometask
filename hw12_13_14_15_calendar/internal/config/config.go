package config

import (
	"github.com/spf13/viper"
)

func Load(filepath string) (c *Config, err error) {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("toml")
	setDefaults()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("log.log_level", "info")
}

type Config struct {
	Log    Log
	Server Server
	App    App
	PSQL   PSQL
}

type PSQL struct {
	DSN string
}

type Log struct {
	LogFile  string `mapstructure:"log_file"`
	LogLevel string `mapstructure:"log_level"`
}

type Server struct {
	Addr         string
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

type App struct {
	RepoType string `mapstructure:"repo_type"`
}
