package parser

import (
	"regexp"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

var hunkHeaderPattern = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)

type Parser struct {
	detector FunctionDetector
}

func NewParser() *Parser {
	return &Parser{detector: NewFunctionDetector()}
}

func NewParserWithDetector(detector FunctionDetector) *Parser {
	return &Parser{detector: detector}
}

func (p *Parser) ParseDiff(diff string) []model.DiffChunk {
	sections := splitFileSections(diff)
	var chunks []model.DiffChunk

	for _, section := range sections {
		filePath := extractFilePath(section)
		if filePath == "" {
			continue
		}

		language := DetectLanguage(filePath)
		hunks := parseHunks(section)
		for _, hunk := range hunks {
			funcLines := make([]string, 0, len(hunk.contextLines)+len(hunk.addedLines))
			funcLines = append(funcLines, hunk.contextLines...)
			funcLines = append(funcLines, hunk.addedLines...)

			chunk := model.DiffChunk{
				FilePath:     filePath,
				Language:     language,
				Function:     p.detector.DetectFunction(language, funcLines),
				AddedLines:   append([]string(nil), hunk.addedLines...),
				DeletedLines: append([]string(nil), hunk.deletedLines...),
				StartLine:    hunk.newStart,
				EndLine:      hunk.newEnd,
			}
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}

func splitFileSections(diff string) []string {
	lines := strings.Split(diff, "\n")
	var sections []string
	var current []string

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git ") && len(current) > 0 {
			sections = append(sections, strings.Join(current, "\n"))
			current = []string{line}
			continue
		}
		current = append(current, line)
	}

	if len(current) > 0 {
		sections = append(sections, strings.Join(current, "\n"))
	}

	return sections
}

func extractFilePath(section string) string {
	lines := strings.Split(section, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "+++ ") {
			path := strings.TrimPrefix(line, "+++ ")
			path = strings.TrimSpace(path)
			if path == "/dev/null" {
				return ""
			}
			if strings.HasPrefix(path, "b/") {
				path = strings.TrimPrefix(path, "b/")
			}
			return path
		}
	}
	return ""
}

type hunkData struct {
	newStart      int
	newEnd        int
	addedLines    []string
	deletedLines  []string
	contextLines  []string
}

func parseHunks(section string) []hunkData {
	lines := strings.Split(section, "\n")
	var hunks []hunkData
	var current *hunkData
	newLine := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "@@") {
			if current != nil {
				hunks = append(hunks, *current)
			}
			matches := hunkHeaderPattern.FindStringSubmatch(line)
			if len(matches) < 4 {
				current = nil
				continue
			}
			newStart := atoi(matches[3])
			newLine = newStart
			current = &hunkData{newStart: newStart, newEnd: newStart}
			continue
		}

		if current == nil || len(line) == 0 {
			continue
		}

		if line == `\ No newline at end of file` {
			continue
		}

		prefix := line[0]
		content := ""
		if len(line) > 1 {
			content = line[1:]
		}

		switch prefix {
		case '+':
			current.addedLines = append(current.addedLines, content)
			current.newEnd = newLine
			newLine++
		case '-':
			current.deletedLines = append(current.deletedLines, content)
		case ' ':
			current.contextLines = append(current.contextLines, content)
			current.newEnd = newLine
			newLine++
		}
	}

	if current != nil {
		hunks = append(hunks, *current)
	}

	return hunks
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + int(c-'0')
	}
	return n
}
