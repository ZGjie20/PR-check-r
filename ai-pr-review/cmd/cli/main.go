package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/config"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/ai"
	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/repository"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/review"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/service"
)

const configPath = "config/config.yaml"

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	db, err := repository.Connect(&cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := repository.InitSchema(ctx, db); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing database schema: %v\n", err)
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

	ghClient := ghclient.NewClient(cfg.GitHubToken)
	llm := ai.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.APIBase, cfg.Model)
	reviewSvc := service.NewReviewService(
		ghClient,
		parser.NewParser(),
		review.NewEngine(llm),
		repository.NewReviewRepository(db),
		cfg.OutputDir,
	)

	result, err := reviewSvc.CreateReview(ctx, prURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating review: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Review completed.")
	fmt.Println("Result saved to:")
	fmt.Println(result.OutputFile)
}
