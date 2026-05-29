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

type ReviewRecord struct {
	ID        int64         `json:"id"`
	PRTitle   string        `json:"pr_title"`
	PRNumber  int           `json:"pr_number"`
	Issues    []ReviewIssue `json:"issues,omitempty"`
	CreatedAt string        `json:"created_at"`
}

type ReviewListItem struct {
	ID         int64  `json:"id"`
	PRTitle    string `json:"pr_title"`
	PRNumber   int    `json:"pr_number"`
	IssueCount int    `json:"issue_count"`
	CreatedAt  string `json:"created_at"`
}

type ReviewListResult struct {
	Items []ReviewListItem `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type CreateReviewResult struct {
	ID         int64         `json:"id"`
	PRTitle    string        `json:"pr_title"`
	PRNumber   int           `json:"pr_number"`
	Issues     []ReviewIssue `json:"issues"`
	CreatedAt  string        `json:"created_at"`
	OutputFile string        `json:"output_file"`
}
