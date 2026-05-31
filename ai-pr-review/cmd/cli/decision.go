package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
)

func handlePRDecision(
	ctx context.Context,
	reader *bufio.Reader,
	ghClient *ghclient.Client,
	owner, repo string,
	number int,
	result *model.CreateReviewResult,
) error {
	printReviewSummary(result)

	if !promptYesNo(reader, "\nApprove this PR? (y/yes): ") {
		return handleReject(ctx, reader, ghClient, owner, repo, number, result)
	}

	approveBody := strings.TrimSpace(result.ReviewResult.PRChangeSummary)
	if err := ghClient.ApprovePR(ctx, owner, repo, number, approveBody); err != nil {
		return fmt.Errorf("approve PR: %w", err)
	}
	fmt.Println("PR approved on GitHub.")

	if !promptYesNo(reader, "Merge this PR? (y/yes): ") {
		fmt.Println("Approved without merge.")
		return nil
	}

	if err := ghClient.MergePR(ctx, owner, repo, number); err != nil {
		return fmt.Errorf("merge PR: %w", err)
	}
	fmt.Println("PR merged successfully.")
	return nil
}

func handleReject(
	ctx context.Context,
	reader *bufio.Reader,
	ghClient *ghclient.Client,
	owner, repo string,
	number int,
	result *model.CreateReviewResult,
) error {
	draft := ghclient.FormatRejectComment(result)
	comment, err := readMultilineComment(reader, draft)
	if err != nil {
		return fmt.Errorf("read reject comment: %w", err)
	}

	if err := ghClient.RequestChangesPR(ctx, owner, repo, number, comment); err != nil {
		return fmt.Errorf("request changes: %w", err)
	}
	fmt.Println("Requested changes on GitHub.")

	if err := ghClient.CreatePRComment(ctx, owner, repo, number, comment); err != nil {
		return fmt.Errorf("create PR comment: %w", err)
	}
	fmt.Println("Reject comment posted on GitHub.")
	return nil
}

func printReviewSummary(result *model.CreateReviewResult) {
	fmt.Println("\nReview summary:")
	fmt.Printf("  Title:  %s (#%d)\n", result.PRTitle, result.PRNumber)
	fmt.Printf("  Issues: total=%d high=%d medium=%d low=%d\n",
		result.TotalIssues, result.HighIssues, result.MediumIssues, result.LowIssues)

	if summary := strings.TrimSpace(result.ReviewResult.PRChangeSummary); summary != "" {
		fmt.Println("\nChange summary:")
		fmt.Println(summary)
	}
}

func reportPRDecisionError(err error) {
	fmt.Fprintf(os.Stderr, "Error handling PR decision: %v\n", err)
	os.Exit(1)
}
