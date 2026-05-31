package langchain

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/prompt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Reviewer struct {
	llm             llms.Model
	reviewRenderer  *prompt.Renderer
	summaryRenderer *prompt.SummaryRenderer
}

func NewReviewer(apiKey, apiBase, model string, reviewRenderer *prompt.Renderer, summaryRenderer *prompt.SummaryRenderer) (*Reviewer, error) {
	llm, err := openai.New(
		openai.WithToken(apiKey),
		openai.WithBaseURL(apiBase),
		openai.WithModel(model),
		openai.WithResponseFormat(openai.ResponseFormatJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("create langchaingo openai client: %w", err)
	}

	return &Reviewer{
		llm:             llm,
		reviewRenderer:  reviewRenderer,
		summaryRenderer: summaryRenderer,
	}, nil
}

func newReviewerWithModel(llm llms.Model, reviewRenderer *prompt.Renderer, summaryRenderer *prompt.SummaryRenderer) *Reviewer {
	return &Reviewer{
		llm:             llm,
		reviewRenderer:  reviewRenderer,
		summaryRenderer: summaryRenderer,
	}
}

type aiResponseLegacy struct {
	Issues []model.ReviewIssue `json:"issues"`
}

type reviewResultBySeverity struct {
	High   []model.ReviewIssue `json:"high"`
	Medium []model.ReviewIssue `json:"medium"`
	Low    []model.ReviewIssue `json:"low"`
}

type aiResponseV2 struct {
	ReviewResult reviewResultBySeverity `json:"review_result"`
}

func (r *Reviewer) ReviewCode(ctx context.Context, input model.ReviewInput) (*model.ReviewResult, error) {
	userPrompt, err := r.reviewRenderer.RenderUser(input)
	if err != nil {
		return nil, err
	}

	resp, err := r.llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: r.reviewRenderer.SystemMessage()}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: userPrompt}},
		},
	}, llms.WithTemperature(0.2))
	if err != nil {
		return nil, fmt.Errorf("chat completion: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return nil, fmt.Errorf("llm returned no choices")
	}

	return parseReviewResponse(resp.Choices[0].Content, input.Chunk.FilePath)
}

type summaryResponse struct {
	PRChangeSummary string `json:"pr_change_summary"`
}

func (r *Reviewer) SummarizePR(ctx context.Context, input model.SummaryInput) (string, error) {
	userPrompt, err := r.summaryRenderer.RenderUser(input)
	if err != nil {
		return "", err
	}

	resp, err := r.llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: r.summaryRenderer.SystemMessage()}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: userPrompt}},
		},
	}, llms.WithTemperature(0.2))
	if err != nil {
		return "", fmt.Errorf("chat completion for summary: %w", err)
	}

	if resp == nil || len(resp.Choices) == 0 {
		return "", fmt.Errorf("llm returned no choices for summary")
	}

	return parseSummaryResponse(resp.Choices[0].Content)
}

func parseSummaryResponse(content string) (string, error) {
	content = strings.TrimSpace(content)

	var parsed summaryResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		snippet := content
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return "", fmt.Errorf("parse summary response: %w (content: %s)", err, snippet)
	}
	return strings.TrimSpace(parsed.PRChangeSummary), nil
}

func parseReviewResponse(content, defaultFile string) (*model.ReviewResult, error) {
	content = strings.TrimSpace(content)

	var raw map[string]json.RawMessage
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		snippet := content
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return nil, fmt.Errorf("parse ai response: %w (content: %s)", err, snippet)
	}

	if _, ok := raw["review_result"]; ok {
		var parsed aiResponseV2
		if err := json.Unmarshal([]byte(content), &parsed); err != nil {
			return nil, fmt.Errorf("parse ai response v2: %w", err)
		}
		issues := flattenReviewResult(parsed.ReviewResult, defaultFile)
		return &model.ReviewResult{Issues: issues}, nil
	}

	var parsed aiResponseLegacy
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		return nil, fmt.Errorf("parse ai response legacy: %w", err)
	}
	fillDefaultFile(parsed.Issues, defaultFile)
	return &model.ReviewResult{Issues: parsed.Issues}, nil
}

func flattenReviewResult(result reviewResultBySeverity, defaultFile string) []model.ReviewIssue {
	issues := make([]model.ReviewIssue, 0,
		len(result.High)+len(result.Medium)+len(result.Low))

	for _, issue := range result.High {
		issue.Severity = "high"
		issues = append(issues, issue)
	}
	for _, issue := range result.Medium {
		issue.Severity = "medium"
		issues = append(issues, issue)
	}
	for _, issue := range result.Low {
		issue.Severity = "low"
		issues = append(issues, issue)
	}

	fillDefaultFile(issues, defaultFile)
	return issues
}

func fillDefaultFile(issues []model.ReviewIssue, defaultFile string) {
	for i := range issues {
		if issues[i].File == "" {
			issues[i].File = defaultFile
		}
	}
}
