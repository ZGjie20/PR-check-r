package parser

import (
	"regexp"
	"strings"
)

var goFuncPattern = regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

type FunctionDetector interface {
	DetectFunction(language string, lines []string) string
}

type defaultFunctionDetector struct{}

func NewFunctionDetector() FunctionDetector {
	return &defaultFunctionDetector{}
}

func (d *defaultFunctionDetector) DetectFunction(language string, lines []string) string {
	switch language {
	case "go":
		return detectGoFunction(lines)
	default:
		return ""
	}
}

func detectGoFunction(lines []string) string {
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		matches := goFuncPattern.FindStringSubmatch(line)
		if len(matches) == 2 {
			return matches[1]
		}
	}
	return ""
}
