package main

import (
	"github.com/lablabs/aws-service-quotas-exporter/internal/app"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/flags"
	log "github.com/lablabs/aws-service-quotas-exporter/pkg/log"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/service"
	"os"
)

func main() {
	cfg := app.Config{}
	flags.ParseOrFail(&cfg, os.Args)
	logger := log.NewLoggerOrFail(cfg.Log.Format, cfg.Log.Level)
	app, err := app.NewApplication(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	ctx := service.SignContext()
	if err := app.Run(ctx); err != nil {
		logger.Fatal(err)
	}
}
