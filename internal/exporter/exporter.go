package exporter

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

type Collector interface {
	Register(r *prometheus.Registry) error
	Collect(g *errgroup.Group, ctx context.Context)
}

const (
	defaultScrapeInterval   = time.Second * 60
	defaultCollectorTimeout = time.Second * 5
)

func NewExporter(log *logrus.Logger, cls []Collector, options ...Option) (*Exporter, error) {
	cfg := config{
		interval: defaultScrapeInterval,
		timeout:  defaultCollectorTimeout,
	}
	for _, o := range options {
		if err := o(&cfg); err != nil {
			return nil, fmt.Errorf("unable to configure exporter: %w", err)
		}
	}
	e := Exporter{
		log: log,
		cfg: &cfg,
		cls: cls,
	}
	return &e, nil
}

type Exporter struct {
	log *logrus.Logger
	cfg *config
	cls []Collector
}

func (e *Exporter) Run(ctx context.Context) error {
	err := e.scrape(ctx)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(e.cfg.interval)
	e.log.Debugf("scrape metrics every: %v", e.cfg.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := e.scrape(ctx)
			if err != nil {
				e.log.Errorf("unable to scrape metric: %v", err)
			}
		}
	}
	return nil
}

func (e *Exporter) scrape(ctx context.Context) error {
	e.log.Debugf("start scraping metrics %v with timeout: %v", time.Now().Format(time.RFC3339), e.cfg.timeout)
	ctx, cancel := context.WithTimeout(ctx, e.cfg.timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, c := range e.cls {
		c.Collect(g, ctx)
	}
	err := g.Wait()
	return err
}
