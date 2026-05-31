package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/output"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/repository"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/review"
)

var ErrReviewNotFound = errors.New("review not found")
var ErrCommentRequired = errors.New("comment is required")

type ReviewService struct {
	ghClient  *ghclient.Client
	parser    *parser.Parser
	engine    *review.Engine
	repo      *repository.ReviewRepository
	outputDir string
	modelName string
}

func NewReviewService(
	ghClient *ghclient.Client,
	diffParser *parser.Parser,
	engine *review.Engine,
	repo *repository.ReviewRepository,
	outputDir string,
	modelName string,
) *ReviewService {
	return &ReviewService{
		ghClient:  ghClient,
		parser:    diffParser,
		engine:    engine,
		repo:      repo,
		outputDir: outputDir,
		modelName: modelName,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, prURL string) (*model.CreateReviewResult, error) {
	owner, repo, number, err := ghclient.ParsePRURL(prURL)
	if err != nil {
		return nil, fmt.Errorf("parse PR URL: %w", err)
	}

	fetchResult, err := s.ghClient.FetchPR(ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("fetch PR: %w", err)
	}

	chunks := s.parser.ParseDiff(fetchResult.PR.Diff)

	reviewResult, err := s.engine.Run(ctx, fetchResult.PR, chunks, fetchResult.Commits)
	if err != nil {
		return nil, fmt.Errorf("run review: %w", err)
	}

	outputPath, err := output.WriteReview(s.outputDir, reviewResult)
	if err != nil {
		return nil, fmt.Errorf("write output: %w", err)
	}

	id, err := s.repo.Save(ctx, &model.ReviewSaveInput{
		Result:       reviewResult,
		RepoName:     owner + "/" + repo,
		PRURL:        prURL,
		AIModel:      s.modelName,
		ReviewStatus: "completed",
	})
	if err != nil {
		return nil, fmt.Errorf("save review: %w", err)
	}

	return &model.CreateReviewResult{
		ID:           id,
		PRTitle:      reviewResult.PRTitle,
		PRNumber:     reviewResult.PRNumber,
		RepoName:     owner + "/" + repo,
		PRURL:        prURL,
		AIModel:      s.modelName,
		ReviewStatus: "completed",
		TotalIssues:  reviewResult.TotalIssues,
		HighIssues:   reviewResult.HighIssues,
		MediumIssues: reviewResult.MediumIssues,
		LowIssues:    reviewResult.LowIssues,
		ReviewResult: reviewResult.ReviewResult,
		RawDiff:      reviewResult.RawDiff,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		OutputFile:   outputPath,
	}, nil
}

func (s *ReviewService) GetReview(ctx context.Context, id int64) (*model.ReviewRecord, error) {
	record, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, ErrReviewNotFound
	}
	return record, nil
}

func (s *ReviewService) ListReviews(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error) {
	items, total, err := s.repo.List(ctx, page, limit, prNumber)
	if err != nil {
		return nil, err
	}
	return &model.ReviewListResult{
		Items: items,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
