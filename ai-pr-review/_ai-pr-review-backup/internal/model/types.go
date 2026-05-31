package model

type PullRequest struct {
	Repo   string
	Number int
	Title  string
	Author string
	Files  []string
	Diff   string
}

type DiffChunk struct {
	FilePath string
	Language string
	Function string

	AddedLines   []string
	DeletedLines []string

	StartLine int
	EndLine   int
}

type ReviewIssue struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Severity   string `json:"severity"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

type AIReviewResult struct {
	PRTitle  string        `json:"pr_title"`
	PRNumber int           `json:"pr_number"`
	Issues   []ReviewIssue `json:"issues"`
}

type ReviewInput struct {
	PRTitle  string
	PRNumber int
	Commits  []string
	Chunk    DiffChunk
}

type ReviewResult struct {
	Issues []ReviewIssue
}
