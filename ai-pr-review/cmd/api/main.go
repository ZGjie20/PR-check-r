package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/api"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/config"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/ai/langchain"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/handler"
	ghclient "github.com/ZGjie20/PR-check-r/ai-pr-review/internal/github"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/parser"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/prompt"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/repository"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/review"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/service"
	"github.com/gin-gonic/gin"
)

const configPath = "config/config.yaml"

func main() {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx := context.Background()

	db, err := repository.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := repository.InitSchema(ctx, db); err != nil {
		log.Fatalf("init schema: %v", err)
	}

	ghClient := ghclient.NewClient(cfg.GitHubToken)

	promptTemplates, err := prompt.Load(cfg.PromptDir)
	if err != nil {
		log.Fatalf("load prompts: %v", err)
	}
	summaryTemplates, err := prompt.LoadSummary(cfg.PromptDir)
	if err != nil {
		log.Fatalf("load summary prompts: %v", err)
	}
	promptRenderer, err := prompt.NewRenderer(promptTemplates)
	if err != nil {
		log.Fatalf("init prompt renderer: %v", err)
	}
	summaryRenderer, err := prompt.NewSummaryRenderer(summaryTemplates)
	if err != nil {
		log.Fatalf("init summary prompt renderer: %v", err)
	}

	llm, err := langchain.NewReviewer(cfg.OpenAIAPIKey, cfg.APIBase, cfg.Model, promptRenderer, summaryRenderer)
	if err != nil {
		log.Fatalf("init llm: %v", err)
	}
	reviewSvc := service.NewReviewService(
		ghClient,
		parser.NewParser(),
		review.NewEngine(llm),
		repository.NewReviewRepository(db),
		cfg.OutputDir,
		cfg.Model,
	)

	reviewHandler := handler.NewReviewHandler(reviewSvc)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api.RegisterRoutes(router, reviewHandler)

	addr := cfg.ServerAddr()
	srvErr := make(chan error, 1)
	go func() {
		log.Printf("API server listening on %s", addr)
		srvErr <- router.Run(addr)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		log.Fatalf("server error: %v", err)
	case sig := <-quit:
		fmt.Printf("shutting down: %v\n", sig)
	}
}
