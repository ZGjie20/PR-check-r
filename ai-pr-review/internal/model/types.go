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
	RawDiff      string

	StartLine int
	EndLine   int
}

type ReviewIssue struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Severity   string `json:"severity,omitempty"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

type ReviewIssueDetail struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
}

func (i ReviewIssue) ToDetail() ReviewIssueDetail {
	return ReviewIssueDetail{
		File:       i.File,
		Line:       i.Line,
		Message:    i.Message,
		Suggestion: i.Suggestion,
	}
}

type ReviewResultBySeverity struct {
	PRChangeSummary string              `json:"pr_change_summary"`
	High            []ReviewIssueDetail `json:"high"`
	Medium          []ReviewIssueDetail `json:"medium"`
	Low             []ReviewIssueDetail `json:"low"`
}

func GroupIssuesBySeverity(issues []ReviewIssue) ReviewResultBySeverity {
	grouped := ReviewResultBySeverity{
		High:   []ReviewIssueDetail{},
		Medium: []ReviewIssueDetail{},
		Low:    []ReviewIssueDetail{},
	}
	for _, issue := range issues {
		detail := issue.ToDetail()
		switch issue.Severity {
		case "high":
			grouped.High = append(grouped.High, detail)
		case "medium":
			grouped.Medium = append(grouped.Medium, detail)
		case "low":
			grouped.Low = append(grouped.Low, detail)
		default:
			grouped.Medium = append(grouped.Medium, detail)
		}
	}
	return grouped
}

func CountIssues(grouped ReviewResultBySeverity) (total, high, medium, low int) {
	high = len(grouped.High)
	medium = len(grouped.Medium)
	low = len(grouped.Low)
	total = high + medium + low
	return total, high, medium, low
}

type AIReviewResult struct {
	PRTitle      string                 `json:"pr_title"`
	PRNumber     int                    `json:"pr_number"`
	TotalIssues  int                    `json:"total_issues"`
	HighIssues   int                    `json:"high_issues"`
	MediumIssues int                    `json:"medium_issues"`
	LowIssues    int                    `json:"low_issues"`
	ReviewResult ReviewResultBySeverity `json:"review_result"`
	RawDiff      string                 `json:"raw_diff"`
	Issues       []ReviewIssue          `json:"-"`
}

type ReviewInput struct {
	PRTitle  string
	PRNumber int
	Commits  []string
	Chunk    DiffChunk
}

type SummaryInput struct {
	PRTitle  string
	PRNumber int
	Author   string
	Files    []string
	Commits  []string
	RawDiff  string
}

type ReviewResult struct {
	Issues []ReviewIssue
}

type ReviewRecord struct {
	ID           int64                  `json:"id"`
	PRTitle      string                 `json:"pr_title"`
	PRNumber     int                    `json:"pr_number"`
	RepoName     string                 `json:"repo_name,omitempty"`
	PRURL        string                 `json:"pr_url"`
	AIModel      string                 `json:"ai_model,omitempty"`
	ReviewStatus string                 `json:"review_status"`
	TotalIssues  int                    `json:"total_issues"`
	HighIssues   int                    `json:"high_issues"`
	MediumIssues int                    `json:"medium_issues"`
	LowIssues    int                    `json:"low_issues"`
	ReviewResult ReviewResultBySeverity `json:"review_result"`
	RawDiff      string                 `json:"raw_diff"`
	CreatedAt    string                 `json:"created_at"`
}

type ReviewListItem struct {
	ID          int64  `json:"id"`
	PRTitle     string `json:"pr_title"`
	PRNumber    int    `json:"pr_number"`
	TotalIssues int    `json:"total_issues"`
	CreatedAt   string `json:"created_at"`
}

type ReviewListResult struct {
	Items []ReviewListItem `json:"items"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Limit int              `json:"limit"`
}

type CreateReviewResult struct {
	ID           int64                  `json:"id"`
	PRTitle      string                 `json:"pr_title"`
	PRNumber     int                    `json:"pr_number"`
	RepoName     string                 `json:"repo_name,omitempty"`
	PRURL        string                 `json:"pr_url"`
	AIModel      string                 `json:"ai_model,omitempty"`
	ReviewStatus string                 `json:"review_status"`
	TotalIssues  int                    `json:"total_issues"`
	HighIssues   int                    `json:"high_issues"`
	MediumIssues int                    `json:"medium_issues"`
	LowIssues    int                    `json:"low_issues"`
	ReviewResult ReviewResultBySeverity `json:"review_result"`
	RawDiff      string                 `json:"raw_diff"`
	CreatedAt    string                 `json:"created_at"`
	OutputFile   string                 `json:"output_file"`
}

type ReviewSaveInput struct {
	Result       *AIReviewResult
	RepoName     string
	PRURL        string
	AIModel      string
	ReviewStatus string
}

type PRActionResult struct {
	ReviewID int64  `json:"review_id"`
	Action   string `json:"action"`
	Message  string `json:"message"`
}

type RejectCommentDraftResponse struct {
	ReviewID int64  `json:"review_id"`
	Comment  string `json:"comment"`
}
