package app

import (
	"context"
	"fmt"
	"github.com/lablabs/aws-service-quotas-exporter/internal/exporter"
	"github.com/lablabs/aws-service-quotas-exporter/internal/http"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/quotas"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/quota"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

const (
	PrometheusNamespace = "aws_quota_exporter"
)

func NewApplication(log *logrus.Logger, cfg Config) (*Application, error) {
	mng, err := service.NewManager()
	if err != nil {
		return nil, err
	}
	scCfg, err := scrape.LoadAndValidateConfig(cfg.Config)
	if err != nil {
		return nil, fmt.Errorf("unable to configure application: %w", err)
	}

	registry := prometheus.NewRegistry()
	cls := make([]exporter.Collector, 0)
	client, err := quota.NewClient(log)
	if err != nil {
		return nil, err
	}
	qcl, err := quotas.NewCollector(log, scCfg.Quotas, PrometheusNamespace, client)
	if err != nil {
		return nil, err
	}
	cls = append(cls, qcl)

	scl, err := script.NewCollector(log, scCfg.Metrics, PrometheusNamespace)
	if err != nil {
		return nil, err
	}
	cls = append(cls, scl)

	exp, err := exporter.NewExporter(log, cls, registry, exporterOptions(scCfg)...)
	if err != nil {
		return nil, err
	}
	mng.Add(exp)

	http, err := http.NewHTTP(log, cfg.Address, registry)
	if err != nil {
		return nil, err
	}
	mng.Add(http)
	a := Application{
		log: log,
		cfg: cfg,
		mng: mng,
	}
	return &a, nil
}

type Application struct {
	log *logrus.Logger
	cfg Config
	mng *service.Manager
}

func (a *Application) Run(ctx context.Context) error {
	a.log.Infof("exporter is starting")
	err := a.mng.StartAndWait(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()
	a.log.Infof("exporter exit OK")
	return nil
}

func exporterOptions(cfg *scrape.Config) []exporter.Option {
	return []exporter.Option{
		exporter.WithInterval(cfg.Scrape.Interval),
		exporter.WithTimeout(cfg.Scrape.Timeout),
	}
}
