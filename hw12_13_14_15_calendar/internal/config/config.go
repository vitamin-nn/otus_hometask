package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type ServerCfg struct {
	Log        Log
	Server     Server
	GrpcServer GrpcServer
	App        App
	PSQL       PSQL
}

func LoadServer() (c *ServerCfg, err error) {
	cfg := new(ServerCfg)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cs *ServerCfg) Fields() log.Fields {
	return log.Fields{
		"http_server_addr": cs.Server.Addr,
		"grpc_server_addr": cs.GrpcServer.Addr,
		"repo":             cs.App.RepoType,
		"log_level":        cs.Log.LogLevel,
	}
}

type SchedulerCfg struct {
	Log       Log
	App       App
	PSQL      PSQL
	Rabbit    Rabbit
	Scheduler Scheduler
}

func LoadScheduler() (c *SchedulerCfg, err error) {
	cfg := new(SchedulerCfg)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type SenderCfg struct {
	Log    Log
	Rabbit Rabbit
}

func LoadSender() (c *SenderCfg, err error) {
	cfg := new(SenderCfg)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type TestCfg struct {
	App    App
	PSQL   PSQL
	Rabbit Rabbit
	Server Server
}

func LoadTest() (c *TestCfg, err error) {
	cfg := new(TestCfg)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

type PSQL struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DB       string `env:"POSTGRES_DB,required"`
	DBHost   string `env:"POSTGRES_DB_HOST,required"`
	Port     int    `env:"POSTGRES_PORT,required"`
}

func (cpg *PSQL) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cpg.DBHost, cpg.Port, cpg.User, cpg.Password, cpg.DB)
}

type Log struct {
	LogFile  string `env:"CALENDAR_LOG_FILE"`
	LogLevel string `env:"CALENDAR_LOG_LEVEL" envDefault:"info"`
}

type Server struct {
	Addr         string        `env:"CALENDAR_HTTP_ADDR,required"`
	ReadTimeout  time.Duration `env:"CALENDAR_HTTP_READ_TO" envDefault:"15s"`
	WriteTimeout time.Duration `env:"CALENDAR_HTTP_WRITE_TO" envDefault:"15s"`
}

type GrpcServer struct {
	Addr string `env:"CALENDAR_GRPC_ADDR,required"`
}

type App struct {
	RepoType string `env:"CALENDAR_REPO_TYPE" envDefault:"internal"`
}

type Rabbit struct {
	User         string `env:"RABBITMQ_DEFAULT_USER,required"`
	Password     string `env:"RABBITMQ_DEFAULT_PASS,required"`
	Host         string `env:"RABBITMQ_HOST,required"`
	VHost        string `env:"RABBITMQ_DEFAULT_VHOST,required"`
	ExchangeName string `env:"RABBITMQ_EXCHANGE_NAME,required"`
	ExchangeType string `env:"RABBITMQ_EXCHANGE_TYPE,required"`
	QueueName    string `env:"RABBITMQ_QUEUE_NAME,required"`
}

func (cr *Rabbit) GetDSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:5672/%s", cr.User, cr.Password, cr.Host, cr.VHost)
}

type Scheduler struct {
	CheckInterval time.Duration `env:"SCHEDULER_CHECK_INTERVAL" envDefault:"60s"`
	CleanInterval time.Duration `env:"SCHEDULER_CLEAN_INTERVAL" envDefault:"60s"`
	EventLiveDays int           `env:"SCHEDULER_EVENT_LIVE_DAYS" envDefault:"365"`
}
