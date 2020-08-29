package config

import (
	"time"

	log "github.com/sirupsen/logrus"
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
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.read_timeout", "15s")
}

func (cfg *Config) Fields() log.Fields {
	return log.Fields{
		"http_server_addr": cfg.Server.Addr,
		"grpc_server_addr": cfg.GrpcServer.Addr,
		"repo":             cfg.App.RepoType,
		"log_level":        cfg.Log.LogLevel,
	}
}

type Config struct {
	Log        Log
	Server     Server
	GrpcServer GrpcServer `mapstructure:"grpc_server"`
	App        App
	PSQL       PSQL
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
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type GrpcServer struct {
	Addr string
}

type App struct {
	RepoType string `mapstructure:"repo_type"`
}
