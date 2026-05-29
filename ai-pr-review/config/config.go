package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultAPIBase = "https://api.deepseek.com/v1"
const DefaultOutputDir = "output"
const DefaultServerHost = "0.0.0.0"
const DefaultServerPort = 8080

var envRefPattern = regexp.MustCompile(`^\$\{([A-Za-z_][A-Za-z0-9_]*)\}$|^\$([A-Za-z_][A-Za-z0-9_]*)$`)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type databaseFileConfig struct {
	Host     string      `yaml:"host"`
	Port     interface{} `yaml:"port"`
	User     string      `yaml:"user"`
	Password string      `yaml:"password"`
	Name     string      `yaml:"name"`
}

type fileConfig struct {
	GitHubToken  string             `yaml:"github_token"`
	OpenAIAPIKey string             `yaml:"openai_api_key"`
	APIBase      string             `yaml:"api_base"`
	Model        string             `yaml:"model"`
	OutputDir    string             `yaml:"output_dir"`
	Server       ServerConfig       `yaml:"server"`
	Database     databaseFileConfig `yaml:"database"`
}

type Config struct {
	GitHubToken  string         `yaml:"github_token"`
	OpenAIAPIKey string         `yaml:"openai_api_key"`
	APIBase      string         `yaml:"api_base"`
	Model        string         `yaml:"model"`
	OutputDir    string         `yaml:"output_dir"`
	Server       ServerConfig   `yaml:"server"`
	Database     DatabaseConfig `yaml:"database"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file: %w", err)
	}

	var raw fileConfig
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	dbPort, err := resolveConfigPort(raw.Database.Port)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		GitHubToken:  resolveConfigValue(raw.GitHubToken),
		OpenAIAPIKey: resolveConfigValue(raw.OpenAIAPIKey),
		APIBase:      resolveConfigValue(raw.APIBase),
		Model:        resolveConfigValue(raw.Model),
		OutputDir:    resolveConfigValue(raw.OutputDir),
		Server: ServerConfig{
			Host: resolveConfigValue(raw.Server.Host),
			Port: raw.Server.Port,
		},
		Database: DatabaseConfig{
			Host:     resolveConfigValue(raw.Database.Host),
			Port:     dbPort,
			User:     resolveConfigValue(raw.Database.User),
			Password: resolveConfigValue(raw.Database.Password),
			Name:     resolveConfigValue(raw.Database.Name),
		},
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
	if strings.TrimSpace(cfg.OutputDir) == "" {
		cfg.OutputDir = DefaultOutputDir
	}
	if strings.TrimSpace(cfg.Server.Host) == "" {
		cfg.Server.Host = DefaultServerHost
	}
	if cfg.Server.Port <= 0 {
		cfg.Server.Port = DefaultServerPort
	}
	if err := cfg.Database.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (d *DatabaseConfig) validate() error {
	if strings.TrimSpace(d.Host) == "" {
		return fmt.Errorf("database.host is required in config")
	}
	if d.Port <= 0 {
		return fmt.Errorf("database.port is required in config")
	}
	if strings.TrimSpace(d.User) == "" {
		return fmt.Errorf("database.user is required in config")
	}
	if strings.TrimSpace(d.Name) == "" {
		return fmt.Errorf("database.name is required in config")
	}
	return nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name)
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func resolveConfigValue(value string) string {
	value = strings.TrimSpace(value)
	matches := envRefPattern.FindStringSubmatch(value)
	if matches == nil {
		return value
	}

	envName := matches[1]
	if envName == "" {
		envName = matches[2]
	}
	return strings.TrimSpace(os.Getenv(envName))
}

func resolveConfigPort(value interface{}) (int, error) {
	if value == nil {
		return 0, fmt.Errorf("database.port is required in config")
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		resolved := resolveConfigValue(v)
		if resolved == "" {
			return 0, fmt.Errorf("database.port is required in config")
		}
		port, err := strconv.Atoi(resolved)
		if err != nil {
			return 0, fmt.Errorf("database.port must be a number: %w", err)
		}
		return port, nil
	default:
		return 0, fmt.Errorf("database.port must be a number")
	}
}
