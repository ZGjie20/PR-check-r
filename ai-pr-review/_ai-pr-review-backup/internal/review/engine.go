package review

import (
	"context"
	"fmt"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/ai"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

type Engine struct {
	llm ai.LLM
}

func NewEngine(llm ai.LLM) *Engine {
	return &Engine{llm: llm}
}

func (e *Engine) Run(ctx context.Context, pr *model.PullRequest, chunks []model.DiffChunk, commits []string) (*model.AIReviewResult, error) {
	result := &model.AIReviewResult{
		PRTitle:  pr.Title,
		PRNumber: pr.Number,
		Issues:   []model.ReviewIssue{},
	}

	for _, chunk := range chunks {
		if len(chunk.AddedLines) == 0 && len(chunk.DeletedLines) == 0 {
			continue
		}

		input := model.ReviewInput{
			PRTitle:  pr.Title,
			PRNumber: pr.Number,
			Commits:  commits,
			Chunk:    chunk,
		}

		chunkResult, err := e.llm.ReviewCode(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("review chunk %s (%d-%d): %w", chunk.FilePath, chunk.StartLine, chunk.EndLine, err)
		}

		result.Issues = append(result.Issues, chunkResult.Issues...)
	}

	return result, nil
}
