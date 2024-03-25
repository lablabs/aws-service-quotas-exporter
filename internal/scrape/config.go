package scrape

import (
	"fmt"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/quotas"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/config"
	"time"
)

type Scrape struct {
	Interval time.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

func (s *Scrape) Validate() error {
	if s.Interval != 0 && s.Interval < time.Minute {
		return fmt.Errorf("scrape.interval is not valid. Minimal value is 60s")
	}
	if s.Timeout != 0 && s.Timeout < (time.Second*5) {
		return fmt.Errorf("scrape.timeout is not valid. Minimal value is 5s")
	}
	return nil
}

type Config struct {
	Scrape  Scrape          `json:"scrape,omitempty" yaml:"scrape,omitempty"`
	Quotas  []quotas.Config `json:"quotas,omitempty" yaml:"quotas,omitempty"`
	Metrics []script.Config `json:"metrics,omitempty" yaml:"metrics,omitempty"`
}

func (c *Config) Validate() error {
	for _, q := range c.Quotas {
		if err := q.Validate(); err != nil {
			return err
		}
	}
	for _, m := range c.Metrics {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return c.Scrape.Validate()
}

func LoadAndValidateConfig(path string) (*Config, error) {
	cfg := Config{}
	err := config.ParseYamlFromFile(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse metric scrape from file: %w", err)
	}
	err = cfg.Validate()
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type Validator interface {
	Validate() error
}
