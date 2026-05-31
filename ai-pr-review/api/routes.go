package api

import (
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/handler"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, reviewHandler *handler.ReviewHandler) {
	router.Use(middleware.Logger(), middleware.Recovery())

	router.GET("/health", reviewHandler.Health)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/reviews", reviewHandler.Create)
		v1.GET("/reviews", reviewHandler.List)
		v1.GET("/reviews/:id", reviewHandler.Get)
		v1.POST("/reviews/:id/approve", reviewHandler.Approve)
		v1.POST("/reviews/:id/merge", reviewHandler.Merge)
		v1.GET("/reviews/:id/reject-comment-draft", reviewHandler.GetRejectCommentDraft)
		v1.POST("/reviews/:id/reject", reviewHandler.Reject)
	}
}
