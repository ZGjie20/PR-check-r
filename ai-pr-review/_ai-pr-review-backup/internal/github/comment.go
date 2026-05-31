package github

import (
	"fmt"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func FormatRejectComment(result *model.CreateReviewResult) string {
	return FormatRejectCommentFromReview(
		result.ReviewResult,
		result.TotalIssues,
		result.HighIssues,
		result.MediumIssues,
		result.LowIssues,
	)
}

func FormatRejectCommentFromReview(
	reviewResult model.ReviewResultBySeverity,
	total, high, medium, low int,
) string {
	var b strings.Builder

	fmt.Fprintln(&b, "## AI PR Review — 请求修改")
	fmt.Fprintln(&b)

	if summary := strings.TrimSpace(reviewResult.PRChangeSummary); summary != "" {
		fmt.Fprintln(&b, "**变更摘要：**")
		fmt.Fprintln(&b, summary)
		fmt.Fprintln(&b)
	}

	fmt.Fprintf(&b, "**问题统计：** 共 %d 个（高 %d / 中 %d / 低 %d）\n\n", total, high, medium, low)

	appendIssueSection(&b, "高优先级", reviewResult.High)
	appendIssueSection(&b, "中优先级", reviewResult.Medium)
	appendIssueSection(&b, "低优先级", reviewResult.Low)

	return strings.TrimRight(b.String(), "\n")
}

func appendIssueSection(b *strings.Builder, title string, issues []model.ReviewIssueDetail) {
	if len(issues) == 0 {
		return
	}

	fmt.Fprintf(b, "### %s\n", title)
	for _, issue := range issues {
		fmt.Fprintf(b, "- `%s:%d` — %s\n", issue.File, issue.Line, issue.Message)
		if suggestion := strings.TrimSpace(issue.Suggestion); suggestion != "" {
			fmt.Fprintf(b, "  - 建议：%s\n", suggestion)
		}
	}
	fmt.Fprintln(b)
}
