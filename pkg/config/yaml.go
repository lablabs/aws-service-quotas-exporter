package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// ParseYaml bytes to scrape struct
func ParseYaml(data []byte, c interface{}) error {
	err := yaml.Unmarshal(data, c)
	if err != nil {
		return fmt.Errorf("invalid yaml config: %w", err)
	}
	return nil
}

func ParseYamlFromFile(path string, c interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	return ParseYaml(data, c)
}
