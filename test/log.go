package test

import (
	logUtil "github.com/lablabs/aws-service-quotas-exporter/pkg/log"
	log "github.com/sirupsen/logrus"
)

const (
	debugLevel = "DEBUG"
	format     = "json "
)

func DefaultLogger() *log.Logger {
	logger := logUtil.NewLoggerOrFail(format, debugLevel)
	return logger
}
