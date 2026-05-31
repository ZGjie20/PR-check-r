package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/ai"
	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/output"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/review"
)

const configPath = "config/config.yaml"

func main() {
	ctx := context.Background()

	cfg, err := model.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	fmt.Print("Enter GitHub PR URL: ")
	reader := bufio.NewReader(os.Stdin)
	prURL, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	prURL = strings.TrimSpace(prURL)

	owner, repo, number, err := ghclient.ParsePRURL(prURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing PR URL: %v\n", err)
		os.Exit(1)
	}

	client := ghclient.NewClient(cfg.GitHubToken)
	fetchResult, err := client.FetchPR(ctx, owner, repo, number)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching PR: %v\n", err)
		os.Exit(1)
	}

	diffParser := parser.NewParser()
	chunks := diffParser.ParseDiff(fetchResult.PR.Diff)

	llm := ai.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.APIBase, cfg.Model)
	engine := review.NewEngine(llm)

	reviewResult, err := engine.Run(ctx, fetchResult.PR, chunks, fetchResult.Commits)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running review: %v\n", err)
		os.Exit(1)
	}

	outputPath, err := output.WriteReview("output", reviewResult)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Review completed.")
	fmt.Println("Result saved to:")
	fmt.Println(outputPath)
}
