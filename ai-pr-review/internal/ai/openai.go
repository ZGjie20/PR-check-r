package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

type OpenAIClient struct {
	client *openai.Client
	model  string
}

func NewOpenAIClient(apiKey, apiBase, model string) *OpenAIClient {
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = apiBase
	return &OpenAIClient{
		client: openai.NewClientWithConfig(cfg),
		model:  model,
	}
}

type aiResponse struct {
	Issues []model.ReviewIssue `json:"issues"`
}

func (c *OpenAIClient) ReviewCode(ctx context.Context, input model.ReviewInput) (*model.ReviewResult, error) {
	prompt := BuildReviewPrompt(input)

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一名专业的代码审查专家。请用简体中文撰写 message 和 suggestion 字段，且只返回合法 JSON。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		},
		Temperature: 0.2,
	})
	if err != nil {
		return nil, fmt.Errorf("chat completion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("llm returned no choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
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
			parsed.Issues[i].File = input.Chunk.FilePath
		}
	}

	return &model.ReviewResult{Issues: parsed.Issues}, nil
}
