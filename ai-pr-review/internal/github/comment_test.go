package github

import (
	"strings"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func TestFormatRejectCommentFromReview(t *testing.T) {
	reviewResult := model.ReviewResultBySeverity{
		PRChangeSummary: "修复 token 硬编码问题",
		High: []model.ReviewIssueDetail{
			{
				File:       "config.yaml",
				Line:       39,
				Message:    "敏感凭据硬编码",
				Suggestion: "使用环境变量",
			},
		},
		Medium: []model.ReviewIssueDetail{},
		Low:    []model.ReviewIssueDetail{},
	}

	got := FormatRejectCommentFromReview(reviewResult, 1, 1, 0, 0)

	if !strings.Contains(got, "## AI PR Review — 请求修改") {
		t.Fatalf("expected header, got:\n%s", got)
	}
	if !strings.Contains(got, "修复 token 硬编码问题") {
		t.Fatalf("expected summary, got:\n%s", got)
	}
	if !strings.Contains(got, "共 1 个（高 1 / 中 0 / 低 0）") {
		t.Fatalf("expected issue stats, got:\n%s", got)
	}
	if !strings.Contains(got, "`config.yaml:39` — 敏感凭据硬编码") {
		t.Fatalf("expected high issue, got:\n%s", got)
	}
	if !strings.Contains(got, "建议：使用环境变量") {
		t.Fatalf("expected suggestion, got:\n%s", got)
	}
	if strings.Contains(got, "### 中优先级") || strings.Contains(got, "### 低优先级") {
		t.Fatalf("expected empty sections to be omitted, got:\n%s", got)
	}
}

func TestFormatRejectCommentFromReviewEmptyIssues(t *testing.T) {
	got := FormatRejectCommentFromReview(model.ReviewResultBySeverity{}, 0, 0, 0, 0)

	if !strings.Contains(got, "共 0 个（高 0 / 中 0 / 低 0）") {
		t.Fatalf("expected zero issue stats, got:\n%s", got)
	}
	if strings.Contains(got, "### 高优先级") {
		t.Fatalf("expected no issue sections, got:\n%s", got)
	}
}

func TestFormatRejectComment(t *testing.T) {
	result := &model.CreateReviewResult{
		TotalIssues:  2,
		HighIssues:   1,
		MediumIssues: 1,
		LowIssues:    0,
		ReviewResult: model.ReviewResultBySeverity{
			PRChangeSummary: "summary",
			High: []model.ReviewIssueDetail{
				{File: "a.go", Line: 1, Message: "high issue"},
			},
			Medium: []model.ReviewIssueDetail{
				{File: "b.go", Line: 2, Message: "medium issue"},
			},
		},
	}

	got := FormatRejectComment(result)
	if !strings.Contains(got, "summary") {
		t.Fatalf("expected summary from CreateReviewResult, got:\n%s", got)
	}
	if !strings.Contains(got, "共 2 个（高 1 / 中 1 / 低 0）") {
		t.Fatalf("expected counts from CreateReviewResult, got:\n%s", got)
	}
}

func TestFormatRejectCommentFromRecord(t *testing.T) {
	record := &model.ReviewRecord{
		TotalIssues:  1,
		HighIssues:   1,
		MediumIssues: 0,
		LowIssues:    0,
		ReviewResult: model.ReviewResultBySeverity{
			PRChangeSummary: "record summary",
			High: []model.ReviewIssueDetail{
				{File: "main.go", Line: 10, Message: "issue"},
			},
		},
	}

	got := FormatRejectCommentFromRecord(record)
	if !strings.Contains(got, "record summary") {
		t.Fatalf("expected summary from ReviewRecord, got:\n%s", got)
	}
	if !strings.Contains(got, "`main.go:10` — issue") {
		t.Fatalf("expected issue from ReviewRecord, got:\n%s", got)
	}
}
