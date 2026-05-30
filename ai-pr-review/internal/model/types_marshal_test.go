package model_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func TestCreateReviewResultJSONShape(t *testing.T) {
	result := model.CreateReviewResult{
		ID:           1,
		PRTitle:      "config+main",
		PRNumber:     1,
		PRURL:        "https://github.com/org/repo/pull/1",
		ReviewStatus: "completed",
		TotalIssues:  1,
		HighIssues:   1,
		ReviewResult: model.ReviewResultBySeverity{
			PRChangeSummary: "新增 a.go",
			High:            []model.ReviewIssueDetail{{File: "a.go", Line: 1, Message: "m", Suggestion: "s"}},
			Medium:          []model.ReviewIssueDetail{},
			Low:             []model.ReviewIssueDetail{},
		},
		RawDiff:    "@@",
		CreatedAt:  "2026-05-29T16:00:00Z",
		OutputFile: "output/a.json",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	content := string(data)
	if strings.Contains(content, `"issues"`) {
		t.Fatalf("CreateReviewResult should not serialize issues: %s", content)
	}
	for _, key := range []string{`"total_issues"`, `"review_result"`, `"pr_change_summary"`, `"raw_diff"`, `"pr_url"`} {
		if !strings.Contains(content, key) {
			t.Errorf("missing %s in %s", key, content)
		}
	}
}

func TestReviewListItemUsesTotalIssues(t *testing.T) {
	item := model.ReviewListItem{
		ID:          1,
		PRTitle:     "test",
		PRNumber:    2,
		TotalIssues: 3,
		CreatedAt:   "2026-05-29T16:00:00Z",
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if !strings.Contains(string(data), `"total_issues":3`) {
		t.Fatalf("expected total_issues field, got %s", data)
	}
	if strings.Contains(string(data), `"issue_count"`) {
		t.Fatalf("should not contain issue_count, got %s", data)
	}
}
