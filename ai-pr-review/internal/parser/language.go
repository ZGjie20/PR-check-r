package parser

import (
	"path/filepath"
	"strings"
)

var extensionLanguage = map[string]string{
	".go":   "go",
	".py":   "python",
	".js":   "javascript",
	".ts":   "typescript",
	".tsx":  "typescript",
	".jsx":  "javascript",
	".java": "java",
	".rs":   "rust",
	".rb":   "ruby",
	".php":  "php",
	".cs":   "csharp",
	".cpp":  "cpp",
	".c":    "c",
	".h":    "c",
	".sql":  "sql",
	".yaml": "yaml",
	".yml":  "yaml",
	".json": "json",
	".md":   "markdown",
	".sh":   "shell",
}

func DetectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	if lang, ok := extensionLanguage[ext]; ok {
		return lang
	}
	return "unknown"
}
