package scrape

import (
	"fmt"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/quotas"
	"github.com/lablabs/aws-service-quotas-exporter/internal/scrape/script"
	"github.com/lablabs/aws-service-quotas-exporter/pkg/config"
	"time"
)

type Global struct {
	Interval time.Duration `json:"interval,omitempty" yaml:"interval,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
}

type Config struct {
	Global  Global          `json:"global,omitempty" yaml:"global,omitempty"`
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
	return nil
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
