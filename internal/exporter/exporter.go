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
	Register(ctx context.Context, r *prometheus.Registry) error
	Collect(ctx context.Context) error
}

const (
	defaultScrapeInterval   = time.Second * 60
	defaultCollectorTimeout = time.Second * 5
)

func NewExporter(log *logrus.Logger, cls []Collector, r *prometheus.Registry, options ...Option) (*Exporter, error) {
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
		r:   r,
	}
	return &e, nil
}

type Exporter struct {
	log *logrus.Logger
	cfg *config
	cls []Collector
	r   *prometheus.Registry
}

func (e *Exporter) Run(ctx context.Context) error {
	err := e.register(ctx)
	if err != nil {
		return err
	}
	ticker := time.NewTicker(e.cfg.interval)
	e.log.Debugf("scrape metrics every: %v", e.cfg.interval)
	defer ticker.Stop()
end:
	for {
		select {
		case <-ctx.Done():
			break end
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
	e.log.Debugf("scrape metrics")
	ctx, cancel := context.WithTimeout(ctx, e.cfg.timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, c := range e.cls {
		cl := c
		g.Go(func() error {
			return cl.Collect(ctx)
		})
	}
	err := g.Wait()
	return err
}

func (e *Exporter) register(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, e.cfg.timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, c := range e.cls {
		cl := c
		g.Go(func() error {
			return cl.Register(ctx, e.r)
		})
	}
	err := g.Wait()
	return err
}
