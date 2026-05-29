package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/api"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/handler"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/service"
	"github.com/gin-gonic/gin"
)

type mockReviewService struct {
	createFn func(ctx context.Context, prURL string) (*model.CreateReviewResult, error)
	getFn    func(ctx context.Context, id int64) (*model.ReviewRecord, error)
	listFn   func(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error)
}

func (m *mockReviewService) CreateReview(ctx context.Context, prURL string) (*model.CreateReviewResult, error) {
	return m.createFn(ctx, prURL)
}

func (m *mockReviewService) GetReview(ctx context.Context, id int64) (*model.ReviewRecord, error) {
	return m.getFn(ctx, id)
}

func (m *mockReviewService) ListReviews(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error) {
	return m.listFn(ctx, page, limit, prNumber)
}

func setupTestRouter(svc *mockReviewService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api.RegisterRoutes(router, handler.NewReviewHandler(svc))
	return router
}

func TestHealth(t *testing.T) {
	router := setupTestRouter(&mockReviewService{})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCreateReviewInvalidURL(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		createFn: func(ctx context.Context, prURL string) (*model.CreateReviewResult, error) {
			t.Fatal("service should not be called for invalid URL")
			return nil, nil
		},
	})

	body := []byte(`{"pr_url":"https://gitlab.com/org/repo/pull/1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestCreateReviewSuccess(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		createFn: func(ctx context.Context, prURL string) (*model.CreateReviewResult, error) {
			return &model.CreateReviewResult{
				ID:         1,
				PRTitle:    "fix login",
				PRNumber:   123,
				Issues:     []model.ReviewIssue{},
				CreatedAt:  "2026-05-29T16:00:00Z",
				OutputFile: "output/fix_login.json",
			}, nil
		},
	})

	body := []byte(`{"pr_url":"https://github.com/org/repo/pull/123"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	var resp model.CreateReviewResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID != 1 || resp.PRNumber != 123 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}

func TestGetReviewNotFound(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		getFn: func(ctx context.Context, id int64) (*model.ReviewRecord, error) {
			return nil, service.ErrReviewNotFound
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}

func TestListReviews(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		listFn: func(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error) {
			return &model.ReviewListResult{
				Items: []model.ReviewListItem{
					{ID: 1, PRTitle: "fix login", PRNumber: 123, IssueCount: 2, CreatedAt: "2026-05-29T16:00:00Z"},
				},
				Total: 1,
				Page:  1,
				Limit: 20,
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews?page=1&limit=20", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp model.ReviewListResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 || len(resp.Items) != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
}
