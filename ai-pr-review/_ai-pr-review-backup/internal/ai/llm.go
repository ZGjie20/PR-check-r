package ai

import (
	"context"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

type LLM interface {
	ReviewCode(ctx context.Context, input model.ReviewInput) (*model.ReviewResult, error)
	SummarizePR(ctx context.Context, input model.SummaryInput) (string, error)
}
