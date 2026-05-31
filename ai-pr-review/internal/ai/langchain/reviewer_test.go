package langchain

import (
	"context"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/prompt"
	"github.com/tmc/langchaingo/llms"
)

type mockLLM struct {
	response string
}

func (m *mockLLM) GenerateContent(_ context.Context, _ []llms.MessageContent, _ ...llms.CallOption) (*llms.ContentResponse, error) {
	return &llms.ContentResponse{
		Choices: []*llms.ContentChoice{{Content: m.response}},
	}, nil
}

func (m *mockLLM) Call(_ context.Context, _ string, _ ...llms.CallOption) (string, error) {
	return m.response, nil
}

func TestReviewCodeParsesResponse(t *testing.T) {
	templates := &prompt.ReviewTemplates{
		System: "system",
		User:   "review {{.PRTitle}}",
	}
	renderer, err := prompt.NewRenderer(templates)
	if err != nil {
		t.Fatalf("NewRenderer() error = %v", err)
	}

	reviewer := newReviewerWithModel(&mockLLM{
		response: `{"issues":[{"file":"","line":12,"severity":"high","message":"问题","suggestion":"建议"}]}`,
	}, renderer)

	result, err := reviewer.ReviewCode(context.Background(), model.ReviewInput{
		PRTitle: "test pr",
		Chunk: model.DiffChunk{
			FilePath: "pkg/main.go",
		},
	})
	if err != nil {
		t.Fatalf("ReviewCode() error = %v", err)
	}

	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].File != "pkg/main.go" {
		t.Errorf("File = %q, want pkg/main.go", result.Issues[0].File)
	}
}

func TestParseReviewResponseV2Format(t *testing.T) {
	content := `{
		"total_issues": 2,
		"high_issues": 1,
		"medium_issues": 0,
		"low_issues": 1,
		"raw_diff": "@@ ...",
		"review_result": {
			"high": [{"file":"","line":42,"message":"高风险","suggestion":"修复高风险"}],
			"medium": [],
			"low": [{"file":"","line":7,"message":"低风险","suggestion":"修复低风险"}]
		}
	}`

	result, err := parseReviewResponse(content, "config.yaml")
	if err != nil {
		t.Fatalf("parseReviewResponse() error = %v", err)
	}
	if len(result.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != "high" || result.Issues[0].File != "config.yaml" {
		t.Errorf("first issue = %+v, want high severity in config.yaml", result.Issues[0])
	}
	if result.Issues[1].Severity != "low" {
		t.Errorf("second issue severity = %q, want low", result.Issues[1].Severity)
	}
}

func TestParseReviewResponseInvalidJSON(t *testing.T) {
	_, err := parseReviewResponse("not json", "pkg/main.go")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
