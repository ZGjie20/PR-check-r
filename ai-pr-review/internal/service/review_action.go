package service

import (
	"context"
	"fmt"
	"strings"

	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func (s *ReviewService) ApproveReview(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
	record, owner, repo, number, err := s.resolvePR(ctx, id)
	if err != nil {
		return nil, err
	}

	body := strings.TrimSpace(comment)
	if body == "" {
		body = strings.TrimSpace(record.ReviewResult.PRChangeSummary)
	}

	if err := s.ghClient.ApprovePR(ctx, owner, repo, number, body); err != nil {
		return nil, fmt.Errorf("approve PR: %w", err)
	}

	return &model.PRActionResult{
		ReviewID: id,
		Action:   "approved",
		Message:  "PR approved on GitHub.",
	}, nil
}

func (s *ReviewService) MergeReview(ctx context.Context, id int64) (*model.PRActionResult, error) {
	_, owner, repo, number, err := s.resolvePR(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.ghClient.MergePR(ctx, owner, repo, number); err != nil {
		return nil, fmt.Errorf("merge PR: %w", err)
	}

	return &model.PRActionResult{
		ReviewID: id,
		Action:   "merged",
		Message:  "PR merged successfully.",
	}, nil
}

func (s *ReviewService) GetRejectCommentDraft(ctx context.Context, id int64) (*model.RejectCommentDraftResponse, error) {
	record, err := s.GetReview(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.RejectCommentDraftResponse{
		ReviewID: id,
		Comment:  ghclient.FormatRejectCommentFromRecord(record),
	}, nil
}

func (s *ReviewService) RejectReview(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
	_, owner, repo, number, err := s.resolvePR(ctx, id)
	if err != nil {
		return nil, err
	}

	body := strings.TrimSpace(comment)
	if body == "" {
		return nil, ErrCommentRequired
	}

	if err := s.ghClient.RequestChangesPR(ctx, owner, repo, number, body); err != nil {
		return nil, fmt.Errorf("request changes: %w", err)
	}

	if err := s.ghClient.CreatePRComment(ctx, owner, repo, number, body); err != nil {
		return nil, fmt.Errorf("create PR comment: %w", err)
	}

	return &model.PRActionResult{
		ReviewID: id,
		Action:   "rejected",
		Message:  "Requested changes and posted comment on GitHub.",
	}, nil
}

func (s *ReviewService) resolvePR(ctx context.Context, id int64) (*model.ReviewRecord, string, string, int, error) {
	record, err := s.GetReview(ctx, id)
	if err != nil {
		return nil, "", "", 0, err
	}

	owner, repo, number, err := ghclient.ParsePRURL(record.PRURL)
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("parse PR URL: %w", err)
	}

	return record, owner, repo, number, nil
}
