package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/service"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/pkg/validator"
	"github.com/gin-gonic/gin"
)

type reviewServicer interface {
	CreateReview(ctx context.Context, prURL string) (*model.CreateReviewResult, error)
	GetReview(ctx context.Context, id int64) (*model.ReviewRecord, error)
	ListReviews(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error)
}

type ReviewHandler struct {
	service reviewServicer
}

func NewReviewHandler(svc reviewServicer) *ReviewHandler {
	return &ReviewHandler{service: svc}
}

type createReviewRequest struct {
	PRURL string `json:"pr_url"`
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var req createReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	prURL, err := validator.ValidatePRURL(req.PRURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateReview(c.Request.Context(), prURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *ReviewHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	record, err := h.service.GetReview(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrReviewNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *ReviewHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	prNumber, _ := strconv.Atoi(c.Query("pr_number"))

	result, err := h.service.ListReviews(c.Request.Context(), page, limit, prNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReviewHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
