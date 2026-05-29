package validator

import (
	"fmt"
	"strings"

	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
)

func ValidatePRURL(rawURL string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("pr_url is required")
	}
	if _, _, _, err := ghclient.ParsePRURL(rawURL); err != nil {
		return "", err
	}
	return rawURL, nil
}
