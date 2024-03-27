package script

import "fmt"

type Config struct {
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Help   string `json:"help,omitempty" yaml:"help,omitempty"`
	Script string `json:"script,omitempty" yaml:"script"`
	Envs   []Env  `json:"env,omitempty" yaml:"envs,omitempty"`
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name of metric is required")
	}
	if c.Help == "" {
		return fmt.Errorf("help of metric is required")
	}
	if c.Script == "" {
		return fmt.Errorf("script is required")
	}
	if c.Envs != nil {
		for _, e := range c.Envs {
			if err := e.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

type Env struct {
	Name  string `json:"name,omitempty" yaml:"name"`
	Value string `json:"value,omitempty" yaml:"value"`
}

func (e *Env) Validate() error {
	if e.Name == "" {
		return fmt.Errorf("name for env is required")
	}
	if e.Value == "" {
		return fmt.Errorf("value for env is required")
	}
	return nil
}

func (c *Config) FormatEnvs() []string {
	evs := make([]string, 0)
	for _, e := range c.Envs {
		evs = append(evs, fmt.Sprintf("%s=%s", e.Name, e.Value))
	}
	return evs
}
