package model

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultAPIBase = "https://api.deepseek.com/v1"

type Config struct {
	GitHubToken  string `yaml:"github_token"`
	OpenAIAPIKey string `yaml:"openai_api_key"`
	APIBase      string `yaml:"api_base"`
	Model        string `yaml:"model"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if strings.TrimSpace(cfg.GitHubToken) == "" {
		return nil, fmt.Errorf("github_token is required in config")
	}
	if strings.TrimSpace(cfg.OpenAIAPIKey) == "" {
		return nil, fmt.Errorf("openai_api_key is required in config")
	}
	if strings.TrimSpace(cfg.Model) == "" {
		return nil, fmt.Errorf("model is required in config")
	}
	if strings.TrimSpace(cfg.APIBase) == "" {
		cfg.APIBase = DefaultAPIBase
	}

	return &cfg, nil
}
