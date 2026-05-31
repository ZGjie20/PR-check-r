package output_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/output"
)

func TestWriteReviewUsesNewJSONFormat(t *testing.T) {
	dir := t.TempDir()
	result := &model.AIReviewResult{
		PRTitle:      "config+main",
		PRNumber:     1,
		TotalIssues:  2,
		HighIssues:   1,
		MediumIssues: 0,
		LowIssues:    1,
		ReviewResult: model.ReviewResultBySeverity{
			High: []model.ReviewIssueDetail{
				{File: "config.yaml", Line: 42, Message: "高风险", Suggestion: "修复"},
			},
			Medium: []model.ReviewIssueDetail{},
			Low: []model.ReviewIssueDetail{
				{File: "main.go", Line: 6, Message: "低风险", Suggestion: "修复"},
			},
		},
		RawDiff: "@@ -1 +1 @@",
		Issues: []model.ReviewIssue{
			{File: "config.yaml", Line: 42, Severity: "high", Message: "高风险", Suggestion: "修复"},
			{File: "main.go", Line: 6, Severity: "low", Message: "低风险", Suggestion: "修复"},
		},
	}

	path, err := output.WriteReview(dir, result)
	if err != nil {
		t.Fatalf("WriteReview() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	content := string(data)
	for _, key := range []string{
		`"pr_title"`,
		`"total_issues"`,
		`"review_result"`,
		`"raw_diff"`,
		`"high"`,
		`"low"`,
	} {
		if !strings.Contains(content, key) {
			t.Errorf("output missing %s\n%s", key, content)
		}
	}
	if strings.Contains(content, `"issues"`) {
		t.Errorf("output should not contain legacy issues field\n%s", content)
	}
	if strings.Contains(content, `"severity"`) {
		t.Errorf("review_result items should not contain severity\n%s", content)
	}

	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if _, ok := decoded["issues"]; ok {
		t.Fatal("decoded output contains issues field")
	}

	if !strings.HasPrefix(filepath.Base(path), "config+main") {
		t.Errorf("unexpected filename: %s", path)
	}
}
