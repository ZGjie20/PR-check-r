package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/config"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/ai/langchain"
	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/prompt"
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

	promptTemplates, err := prompt.Load(cfg.PromptDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading prompts: %v\n", err)
		os.Exit(1)
	}
	promptRenderer, err := prompt.NewRenderer(promptTemplates)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing prompt renderer: %v\n", err)
		os.Exit(1)
	}

	llm, err := langchain.NewReviewer(cfg.OpenAIAPIKey, cfg.APIBase, cfg.Model, promptRenderer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing LLM: %v\n", err)
		os.Exit(1)
	}
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
