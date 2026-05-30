package prompt

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func TestLoadAndRenderSummaryPrompt(t *testing.T) {
	templates, err := LoadSummary(filepath.Join("..", "..", "prompts"))
	if err != nil {
		t.Fatalf("LoadSummary() error = %v", err)
	}

	if strings.TrimSpace(templates.System) == "" {
		t.Fatal("expected non-empty summary system prompt")
	}

	renderer, err := NewSummaryRenderer(templates)
	if err != nil {
		t.Fatalf("NewSummaryRenderer() error = %v", err)
	}

	input := model.SummaryInput{
		PRTitle:  "config+main",
		PRNumber: 1,
		Author:   "alice",
		Files:    []string{"config.yaml", "main.go"},
		Commits:  []string{"add config", "update main"},
		RawDiff:  "diff --git a/config.yaml",
	}

	rendered, err := renderer.RenderUser(input)
	if err != nil {
		t.Fatalf("RenderUser() error = %v", err)
	}

	for _, want := range []string{
		"config+main",
		"PR Number: 1",
		"Author: alice",
		"config.yaml",
		"main.go",
		"add config",
		"Full diff:",
		"diff --git a/config.yaml",
	} {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered summary prompt missing %q\n%s", want, rendered)
		}
	}
}
