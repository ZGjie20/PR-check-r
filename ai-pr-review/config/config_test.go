package config

import (
	"os"
	"testing"
)

func TestResolveConfigValue(t *testing.T) {
	t.Setenv("TEST_TOKEN", "secret-value")

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "literal", input: "deepseek-chat", want: "deepseek-chat"},
		{name: "brace env ref", input: "${TEST_TOKEN}", want: "secret-value"},
		{name: "simple env ref", input: "$TEST_TOKEN", want: "secret-value"},
		{name: "missing env", input: "${NOT_SET_VAR}", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveConfigValue(tt.input); got != tt.want {
				t.Fatalf("resolveConfigValue(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "ghp_test")
	t.Setenv("OPENAI_API_KEY", "sk_test")
	t.Setenv("MODEL", "test-model")
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "3306")
	t.Setenv("DB_USER", "root")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "ai_pr_review")

	dir := t.TempDir()
	path := dir + string(os.PathSeparator) + "config.yaml"
	content := `github_token: ${GITHUB_TOKEN}
openai_api_key: ${OPENAI_API_KEY}
model: ${MODEL}
api_base: "https://example.com/v1"
output_dir: "reviews"
server:
  host: "127.0.0.1"
  port: 9090
database:
  host: ${DB_HOST}
  port: ${DB_PORT}
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  name: ${DB_NAME}
`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GitHubToken != "ghp_test" || cfg.OpenAIAPIKey != "sk_test" || cfg.Model != "test-model" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
	if cfg.APIBase != "https://example.com/v1" {
		t.Fatalf("api_base = %q", cfg.APIBase)
	}
	if cfg.OutputDir != "reviews" {
		t.Fatalf("output_dir = %q", cfg.OutputDir)
	}
	if cfg.Server.Host != "127.0.0.1" || cfg.Server.Port != 9090 {
		t.Fatalf("unexpected server cfg: %+v", cfg.Server)
	}
	if cfg.Database.Host != "localhost" || cfg.Database.Port != 3306 || cfg.Database.Name != "ai_pr_review" {
		t.Fatalf("unexpected database cfg: %+v", cfg.Database)
	}
}
