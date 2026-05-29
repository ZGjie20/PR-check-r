package github_test

import (
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
)

func TestParsePRURL(t *testing.T) {
	cases := []struct {
		url         string
		wantOwner   string
		wantRepo    string
		wantNumber  int
		shouldError bool
	}{
		{
			url:        "https://github.com/org/repo/pull/123",
			wantOwner:  "org",
			wantRepo:   "repo",
			wantNumber: 123,
		},
		{
			url:        "https://github.com/org/repo/pull/123/",
			wantOwner:  "org",
			wantRepo:   "repo",
			wantNumber: 123,
		},
		{
			url:        "http://github.com/my-org/my-repo/pull/42",
			wantOwner:  "my-org",
			wantRepo:   "my-repo",
			wantNumber: 42,
		},
		{
			url:         "https://gitlab.com/org/repo/pull/1",
			shouldError: true,
		},
	}

	for _, tc := range cases {
		owner, repo, number, err := github.ParsePRURL(tc.url)
		if tc.shouldError {
			if err == nil {
				t.Errorf("expected error for URL %s", tc.url)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tc.url, err)
		}
		if owner != tc.wantOwner || repo != tc.wantRepo || number != tc.wantNumber {
			t.Errorf("ParsePRURL(%s) = (%s, %s, %d), want (%s, %s, %d)",
				tc.url, owner, repo, number, tc.wantOwner, tc.wantRepo, tc.wantNumber)
		}
	}
}
