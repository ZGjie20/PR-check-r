package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func WriteReview(dir string, result *model.AIReviewResult) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create output directory: %w", err)
	}

	path := uniqueOutputPath(dir, buildReviewFilename(result))

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal review result: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", fmt.Errorf("write review file: %w", err)
	}

	return path, nil
}

func buildReviewFilename(result *model.AIReviewResult) string {
	title := sanitizeFilename(result.PRTitle)
	if title == "" {
		title = fmt.Sprintf("pr_%d", result.PRNumber)
	}

	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%s.json", title, timestamp)
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"\\", "",
		"/", "",
		":", "",
		"*", "",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "",
	)
	return strings.TrimSpace(replacer.Replace(name))
}

func uniqueOutputPath(dir, filename string) string {
	path := filepath.Join(dir, filename)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return path
	}

	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
