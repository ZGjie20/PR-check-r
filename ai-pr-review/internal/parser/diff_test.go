package parser_test

import (
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
)

const sampleDiff = `diff --git a/service/user.go b/service/user.go
index abc123..def456 100644
--- a/service/user.go
+++ b/service/user.go
@@ -10,6 +10,12 @@ import (
 	"context"
 )
 
+func StartWorker(ctx context.Context) {
+	go func() {
+		for {
+		}
+	}()
+}
+
 func GetUser(id string) (*User, error) {
 	return nil, nil
 }
@@ -20,3 +26,4 @@ func GetUser(id string) (*User, error) {
 func DeleteUser(id string) error {
+	user := GetUserUnsafe(id)
 	return nil
 }
`

func TestParseDiff_HunksAndGoFunction(t *testing.T) {
	p := parser.NewParser()
	chunks := p.ParseDiff(sampleDiff)

	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}

	first := chunks[0]
	if first.FilePath != "service/user.go" {
		t.Errorf("expected file service/user.go, got %s", first.FilePath)
	}
	if first.Language != "go" {
		t.Errorf("expected language go, got %s", first.Language)
	}
	if first.Function != "StartWorker" {
		t.Errorf("expected function StartWorker, got %s", first.Function)
	}
	if first.StartLine != 10 {
		t.Errorf("expected start line 10, got %d", first.StartLine)
	}
	if len(first.AddedLines) != 7 {
		t.Errorf("expected 7 added lines in first hunk, got %d", len(first.AddedLines))
	}

	second := chunks[1]
	if second.Function != "DeleteUser" {
		t.Errorf("expected function DeleteUser, got %s", second.Function)
	}
	if len(second.AddedLines) != 1 {
		t.Errorf("expected 1 added line in second hunk, got %d", len(second.AddedLines))
	}
}

func TestDetectLanguage(t *testing.T) {
	cases := map[string]string{
		"main.go":      "go",
		"script.py":    "python",
		"app.js":       "javascript",
		"unknown.xyz":  "unknown",
	}

	for file, want := range cases {
		got := parser.DetectLanguage(file)
		if got != want {
			t.Errorf("DetectLanguage(%s) = %s, want %s", file, got, want)
		}
	}
}
