package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func promptYesNo(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return isYesAnswer(strings.TrimSpace(line))
}

func isYesAnswer(input string) bool {
	switch strings.ToLower(input) {
	case "y", "yes":
		return true
	default:
		return false
	}
}

func readMultilineComment(reader *bufio.Reader, draft string) (string, error) {
	fmt.Println("\n--- Reject comment draft ---")
	fmt.Println(draft)
	fmt.Println("--- End of draft ---")
	fmt.Print("\nEdit comment? (y to edit, Enter to use draft): ")

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("read edit choice: %w", err)
	}

	if !isYesAnswer(strings.TrimSpace(line)) {
		return draft, nil
	}

	fmt.Println("Enter your comment. End with a single line containing only '.':")
	return readUntilDot(reader)
}

func readUntilDot(reader *bufio.Reader) (string, error) {
	var lines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(lines) > 0 {
				break
			}
			return "", fmt.Errorf("read comment: %w", err)
		}

		if strings.TrimSpace(line) == "." {
			break
		}
		lines = append(lines, strings.TrimRight(line, "\r\n"))
	}

	comment := strings.Join(lines, "\n")
	if strings.TrimSpace(comment) == "" {
		return "", fmt.Errorf("comment cannot be empty")
	}
	return comment, nil
}
