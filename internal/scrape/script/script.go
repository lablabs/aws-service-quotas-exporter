package script

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	for k, _ := range d.Labels {
		r = append(r, k)
	}
	return r
}

func Run(ctx context.Context, cfg Config) ([]Data, error) {
	cmd := exec.CommandContext(ctx, "bash", "-c", cfg.Script)
	envs := make([]string, 0)
	envs = append(envs, os.Environ()...)
	envs = append(envs, cfg.FormatEnvs()...)
	cmd.Env = envs
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout := &bytes.Buffer{}
	cmd.Stdout = stdout
	err = cmd.Run()
	if err != nil {
		errString, err := errorString(stderr)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("script error: %w, std err: %s", err, errString)
	}
	data, err := ParseStdout(stdout)
	if err != nil {
		return nil, fmt.Errorf("unable to parse response from command: %v", cfg.Script)
	}
	return data, nil
}

func errorString(r io.Reader) (string, error) {
	b := strings.Builder{}
	_, err := io.Copy(&b, r)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
