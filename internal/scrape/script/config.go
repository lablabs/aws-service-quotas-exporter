package script

import "fmt"

type Config struct {
	Name    string  `json:"name,omitempty" yaml:"name,omitempty"`
	Help    string  `json:"help,omitempty" yaml:"help,omitempty"`
	Command string  `json:"command,omitempty" yaml:"command,omitempty"`
	Envs    []Env   `json:"env,omitempty" yaml:"envs,omitempty"`
	List    string  `json:"list,omitempty" yaml:"list,omitempty"`
	Value   string  `json:"value,omitempty" yaml:"value,omitempty"`
	Labels  []Label `json:"labels,omitempty" yaml:"labels,omitempty"`
}

type Env struct {
	Name  string `json:"name,omitempty" yaml:"name"`
	Value string `json:"value,omitempty" yaml:"value"`
}

func (c *Config) LabelNames() []string {
	l := make([]string, 0)
	for _, c := range c.Labels {
		l = append(l, c.Name)
	}
	return l
}

func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name attribute is required")
	}
	if c.Command == "" {
		return fmt.Errorf("command attribute is required")
	}
	if c.Value == "" {
		return fmt.Errorf("value attribute is required")
	}
	for _, l := range c.Labels {
		if l.Name == "" {
			return fmt.Errorf("attribute name for label is required")
		}
		if l.GetValue() == "" {
			return fmt.Errorf("jqValue or value for label is required")
		}
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

type Label struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`
	JqValue string `json:"jqValue,omitempty" yaml:"jqValue,omitempty"`
	Value   string `json:"value,omitempty" yaml:"value,omitempty"`
}

func (l Label) GetValue() string {
	if l.JqValue != "" {
		return l.JqValue
	}
	if l.Value != "" {
		return l.Value
	}
	return ""
}
