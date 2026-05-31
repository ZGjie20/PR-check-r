package prompt

import (
	"fmt"
	"os"
	"path/filepath"
)

const reviewSubdir = "review"

type ReviewTemplates struct {
	System string
	User   string
}

func Load(dir string) (*ReviewTemplates, error) {
	if dir == "" {
		dir = "prompts"
	}

	reviewDir := filepath.Join(dir, reviewSubdir)
	systemPath := filepath.Join(reviewDir, "system.txt")
	userPath := filepath.Join(reviewDir, "user.tmpl")

	system, err := os.ReadFile(systemPath)
	if err != nil {
		return nil, fmt.Errorf("read system prompt %s: %w", systemPath, err)
	}

	user, err := os.ReadFile(userPath)
	if err != nil {
		return nil, fmt.Errorf("read user prompt %s: %w", userPath, err)
	}

	return &ReviewTemplates{
		System: string(system),
		User:   string(user),
	}, nil
}
