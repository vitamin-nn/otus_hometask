package logger

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/config"
)

var (
	ErrEmptyLogFile  = errors.New("empty log file error")
	ErrEmptyLogLevel = errors.New("empty log level error")
)

func Init(logCfg config.Log) error {
	err := setLogFile(logCfg.LogFile)
	if err != nil {
		return err
	}

	err = setLogLevel(logCfg.LogLevel)
	if err != nil {
		return err
	}

	return nil
}

func setLogFile(logFile string) error {
	if logFile == "" {
		return ErrEmptyLogFile
	}

	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}
	log.SetOutput(f)

	return nil
}

func setLogLevel(logLevel string) error {
	if logLevel == "" {
		return ErrEmptyLogLevel
	}

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(level)

	return nil
}
