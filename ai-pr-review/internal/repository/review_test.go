package repository_test

import (
	"os"
	"strings"
	"testing"
)

func TestMigrationFilesDefinePRReviewRecord(t *testing.T) {
	data, err := os.ReadFile("../../migrations/002_create_pr_review_record.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}

	required := []string{
		"CREATE TABLE IF NOT EXISTS pr_review_record",
		"total_issues",
		"high_issues",
		"medium_issues",
		"low_issues",
		"review_result",
		"raw_diff",
		"idx_pr_number",
	}
	for _, fragment := range required {
		if !strings.Contains(string(data), fragment) {
			t.Errorf("migration missing %q", fragment)
		}
	}
}
