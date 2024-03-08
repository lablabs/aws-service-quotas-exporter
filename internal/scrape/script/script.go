package script

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ohler55/ojg/jp"
	"os"
	"os/exec"
	"strings"
)

type Data struct {
	Value  float64
	Labels map[string]string
}

func Run(ctx context.Context, cfg Config) ([]Data, error) {

	cs := strings.Split(cfg.Command, " ")
	prg := cs[0]
	args := cs[1:]
	cmd := exec.CommandContext(ctx, prg, args...)
	var stdout, stderr bytes.Buffer
	envs := make([]string, 0)
	envs = append(envs, os.Environ()...)
	envs = append(envs, cfg.FormatEnvs()...)
	cmd.Env = envs
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("script error: %w, std err: %s", err, stderr.String())
	}
	data, err := ParseJSON(stdout.Bytes())
	if err != nil {
		return nil, fmt.Errorf("unable to parse response from command: %v", cfg.Command)
	}

	result := make([]Data, 0)
	if cfg.List != "" {
		items, err := GetArray(data, cfg.List)
		if err != nil {
			return nil, fmt.Errorf("unable to parse list items jq: %w", err)
		}
		for _, it := range items {
			r, err := ParseRecord(it, cfg)
			if err != nil {
				return nil, err
			}
			result = append(result, r)
		}
		return result, nil
	}
	r, err := ParseRecord(data, cfg)
	if err != nil {
		return nil, err
	}
	result = append(result, r)
	return result, nil
}

func ParseRecord(r any, c Config) (Data, error) {
	v, err := GetFloat64(r, c.Value)
	if err != nil {
		return Data{}, err
	}
	labels := make(map[string]string)
	for _, l := range c.Labels {
		lv, err := GetString(r, l.JqValue)
		if err != nil {
			return Data{}, err
		}
		labels[l.Name] = lv
	}
	return Data{
		Value:  v,
		Labels: labels,
	}, nil
}

func ParseJSON(data []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	err := json.Unmarshal(data, &out)
	if err != nil {
		return nil, fmt.Errorf("unable parse to json: %w", err)
	}
	return out, nil
}

func GetFloat64(data any, query string) (float64, error) {
	res, err := ParseJq(data, query)
	if err != nil {
		return 0, err
	}
	v, ok := res[0].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid data. jq expression must return valid float64")
	}
	return v, nil
}

func GetString(data any, query string) (string, error) {
	res, err := ParseJq(data, query)
	if err != nil {
		return "", err
	}
	v, ok := res[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid data. jq expression must return valid string")
	}
	return v, nil
}

func GetArray(data any, query string) ([]any, error) {
	res, err := ParseJq(data, query)
	if err != nil {
		return nil, err
	}
	v, ok := res[0].([]interface{})
	if !ok {
		return v, fmt.Errorf("invalid data. jq expression must return valid string")
	}
	return v, nil
}

func ParseJq(data any, qs string) ([]any, error) {
	q, err := jp.ParseString(qs)
	if err != nil {
		return nil, fmt.Errorf("unable to parse jq selector: %w", err)
	}
	result := q.Get(data)
	if len(result) != 1 {
		return nil, fmt.Errorf("empty data returned from jq expression: %s", qs)
	}
	return result, nil
}
