package script

import (
	"context"
	"fmt"
	"github.com/go-cmd/cmd"
	"os"
)

type Data struct {
	Value  float64
	Labels map[string]string
}

func (d Data) LabelNames() []string {
	if d.Labels == nil {
		return []string{}
	}
	r := make([]string, 0)
	for k := range d.Labels {
		r = append(r, k)
	}
	return r
}

func Run(ctx context.Context, cfg Config) ([]Data, error) {
	c := cmd.NewCmd("bash", "-c", cfg.Script)
	c.Env = append(c.Env, os.Environ()...)
	c.Env = append(c.Env, cfg.FormatEnvs()...)

	<-c.Start()
	select {
	case <-ctx.Done():
		err := c.Stop()
		if err != nil {
			return nil, fmt.Errorf("script error: %w, std err: %s", err, c.Status().Stderr)
		}
	}
	err := c.Status().Error
	if err != nil {
		return nil, fmt.Errorf("script error: %w, std err: %s", err, c.Status().Stderr)
	}
	data, err := ParseStdout(c.Status().Stdout)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response from command: %v", cfg.Script)
	}
	return data, nil
}
