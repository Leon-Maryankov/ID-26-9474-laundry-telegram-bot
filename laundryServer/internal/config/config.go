package config

import (
	"fmt"
	"os"

	yaml "github.com/goccy/go-yaml"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Token string
	Host string `validate:"hostname_port"`
	DatabasePath string `yaml:"databasePath", validate:"filepath"`
}

func NewConfig(path string) (*Config, error) {
	pathContents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config.NewConfig: failed to read file: %w", err)
	}

	result := &Config{}
	if err := yaml.Unmarshal(pathContents, result); err != nil {
		return nil, fmt.Errorf("config.NewConfig: failed to unmarshal file: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(result); err != nil {
		return nil, fmt.Errorf("config.NewConfig: failed to validate config: %w", err)
	}

	return result, nil
}
