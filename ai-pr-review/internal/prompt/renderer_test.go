package prompt

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func TestLoadAndRenderReviewPrompt(t *testing.T) {
	templates, err := Load(filepath.Join("..", "..", "prompts"))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if strings.TrimSpace(templates.System) == "" {
		t.Fatal("expected non-empty system prompt")
	}

	renderer, err := NewRenderer(templates)
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	input := model.ReviewInput{
		PRTitle:  "Fix goroutine leak",
		PRNumber: 42,
		Commits:  []string{"commit one", "commit two"},
		Chunk: model.DiffChunk{
			FilePath:     "internal/foo.go",
			Language:     "go",
			Function:     "Run",
			StartLine:    10,
			EndLine:      20,
			AddedLines:   []string{"go func() {}"},
			DeletedLines: []string{"time.Sleep(1)"},
			RawDiff:      "@@ -10,5 +10,6 @@\n-time.Sleep(1)\n+go func() {}",
		},
	}

	rendered, err := renderer.RenderUser(input)
	if err != nil {
		t.Fatalf("RenderUser() error = %v", err)
	}

	checks := []string{
		"Fix goroutine leak",
		"PR Number: 42",
		"commit one",
		"internal/foo.go",
		"Language: go",
		"Function: Run",
		"Line range: 10-20",
		"+ go func() {}",
		"- time.Sleep(1)",
		"Raw diff:",
		"+go func() {}",
	}
	for _, want := range checks {
		if !strings.Contains(rendered, want) {
			t.Errorf("rendered prompt missing %q\n%s", want, rendered)
		}
	}
}
