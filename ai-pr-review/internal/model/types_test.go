package model

import "testing"

func TestGroupIssuesBySeverity(t *testing.T) {
	issues := []ReviewIssue{
		{File: "a.go", Line: 1, Severity: "high", Message: "h1", Suggestion: "fix h1"},
		{File: "b.go", Line: 2, Severity: "medium", Message: "m1", Suggestion: "fix m1"},
		{File: "c.go", Line: 3, Severity: "low", Message: "l1", Suggestion: "fix l1"},
	}

	grouped := GroupIssuesBySeverity(issues)
	total, high, medium, low := CountIssues(grouped)

	if total != 3 || high != 1 || medium != 1 || low != 1 {
		t.Fatalf("counts = (%d, %d, %d, %d), want (3, 1, 1, 1)", total, high, medium, low)
	}
	if grouped.High[0].File != "a.go" {
		t.Fatalf("high issue detail = %+v", grouped.High[0])
	}
}
