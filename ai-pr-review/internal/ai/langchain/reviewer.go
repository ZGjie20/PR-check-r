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
	llm      llms.Model
	renderer *prompt.Renderer
}

func NewReviewer(apiKey, apiBase, model string, renderer *prompt.Renderer) (*Reviewer, error) {
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
		llm:      llm,
		renderer: renderer,
	}, nil
}

func newReviewerWithModel(llm llms.Model, renderer *prompt.Renderer) *Reviewer {
	return &Reviewer{
		llm:      llm,
		renderer: renderer,
	}
}

type aiResponse struct {
	Issues []model.ReviewIssue `json:"issues"`
}

func (r *Reviewer) ReviewCode(ctx context.Context, input model.ReviewInput) (*model.ReviewResult, error) {
	userPrompt, err := r.renderer.RenderUser(input)
	if err != nil {
		return nil, err
	}

	resp, err := r.llm.GenerateContent(ctx, []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: r.renderer.SystemMessage()}},
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

func parseReviewResponse(content, defaultFile string) (*model.ReviewResult, error) {
	content = strings.TrimSpace(content)
	var parsed aiResponse
	if err := json.Unmarshal([]byte(content), &parsed); err != nil {
		snippet := content
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}
		return nil, fmt.Errorf("parse ai response: %w (content: %s)", err, snippet)
	}

	for i := range parsed.Issues {
		if parsed.Issues[i].File == "" {
			parsed.Issues[i].File = defaultFile
		}
	}

	return &model.ReviewResult{Issues: parsed.Issues}, nil
}
