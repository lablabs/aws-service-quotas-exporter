package quotas

import (
	"fmt"
)

type Config struct {
	ServiceCode string `json:"serviceCode,omitempty" yaml:"serviceCode,omitempty"`
	QuotaCode   string `json:"quotaCode,omitempty" yaml:"quotaCode,omitempty"`
	Region      string `json:"region,omitempty" yaml:"region,omitempty"`
	Default     bool   `json:"default,omitempty" yaml:"default,omitempty"`
}

func (c Config) Validate() error {
	if c.ServiceCode == "" {
		return fmt.Errorf("serviceCode must not be empty")
	}
	if c.QuotaCode == "" {
		return fmt.Errorf("quotaCode must not be empty")
	}
	return nil
}
