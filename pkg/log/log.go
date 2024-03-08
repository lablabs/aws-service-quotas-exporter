package log

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func DefaultLoggerOrFail() *log.Logger {
	logger := NewLoggerOrFail("json", "DEBUG")
	return logger
}

func NewLoggerOrFail(format string, level string) *log.Logger {
	log, err := NewLogger(format, level)
	if err != nil {
		panic(err)
	}
	return log
}

func NewLogger(format string, logLevel string) (*log.Logger, error) {
	logger := log.New()
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	switch format {
	case "json":
		logger.SetFormatter(&log.JSONFormatter{})
	case "text":
		logger.SetFormatter(&log.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}
	logger.SetOutput(os.Stdout)
	logger.SetLevel(level)
	return logger, nil
}
