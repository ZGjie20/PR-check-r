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
	createFn              func(ctx context.Context, prURL string) (*model.CreateReviewResult, error)
	getFn                 func(ctx context.Context, id int64) (*model.ReviewRecord, error)
	listFn                func(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error)
	approveFn             func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error)
	mergeFn               func(ctx context.Context, id int64) (*model.PRActionResult, error)
	rejectDraftFn         func(ctx context.Context, id int64) (*model.RejectCommentDraftResponse, error)
	rejectFn              func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error)
}

func (m *mockReviewService) CreateReview(ctx context.Context, prURL string) (*model.CreateReviewResult, error) {
	if m.createFn != nil {
		return m.createFn(ctx, prURL)
	}
	return nil, nil
}

func (m *mockReviewService) GetReview(ctx context.Context, id int64) (*model.ReviewRecord, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, service.ErrReviewNotFound
}

func (m *mockReviewService) ListReviews(ctx context.Context, page, limit, prNumber int) (*model.ReviewListResult, error) {
	if m.listFn != nil {
		return m.listFn(ctx, page, limit, prNumber)
	}
	return nil, nil
}

func (m *mockReviewService) ApproveReview(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
	if m.approveFn != nil {
		return m.approveFn(ctx, id, comment)
	}
	return nil, nil
}

func (m *mockReviewService) MergeReview(ctx context.Context, id int64) (*model.PRActionResult, error) {
	if m.mergeFn != nil {
		return m.mergeFn(ctx, id)
	}
	return nil, nil
}

func (m *mockReviewService) GetRejectCommentDraft(ctx context.Context, id int64) (*model.RejectCommentDraftResponse, error) {
	if m.rejectDraftFn != nil {
		return m.rejectDraftFn(ctx, id)
	}
	return nil, nil
}

func (m *mockReviewService) RejectReview(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
	if m.rejectFn != nil {
		return m.rejectFn(ctx, id, comment)
	}
	return nil, nil
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
				ID:           1,
				PRTitle:      "fix login",
				PRNumber:     123,
				PRURL:        prURL,
				ReviewStatus: "completed",
				TotalIssues:  1,
				HighIssues:   1,
				ReviewResult: model.ReviewResultBySeverity{
					High: []model.ReviewIssueDetail{
						{File: "service/user.go", Line: 45, Message: "问题", Suggestion: "建议"},
					},
					Medium: []model.ReviewIssueDetail{},
					Low:    []model.ReviewIssueDetail{},
				},
				RawDiff:    "@@ ...",
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
	if resp.ID != 1 || resp.PRNumber != 123 || resp.TotalIssues != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if len(resp.ReviewResult.High) != 1 {
		t.Fatalf("expected 1 high issue, got %+v", resp.ReviewResult)
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
					{ID: 1, PRTitle: "fix login", PRNumber: 123, TotalIssues: 2, CreatedAt: "2026-05-29T16:00:00Z"},
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

func TestApproveReviewSuccess(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		approveFn: func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
			return &model.PRActionResult{
				ReviewID: id,
				Action:   "approved",
				Message:  "PR approved on GitHub.",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews/1/approve", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp model.PRActionResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Action != "approved" {
		t.Fatalf("unexpected action: %+v", resp)
	}
}

func TestMergeReviewSuccess(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		mergeFn: func(ctx context.Context, id int64) (*model.PRActionResult, error) {
			return &model.PRActionResult{
				ReviewID: id,
				Action:   "merged",
				Message:  "PR merged successfully.",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews/1/merge", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp model.PRActionResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Action != "merged" {
		t.Fatalf("unexpected action: %+v", resp)
	}
}

func TestGetRejectCommentDraftSuccess(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		rejectDraftFn: func(ctx context.Context, id int64) (*model.RejectCommentDraftResponse, error) {
			return &model.RejectCommentDraftResponse{
				ReviewID: id,
				Comment:  "## AI PR Review — 请求修改",
			}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews/1/reject-comment-draft", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp model.RejectCommentDraftResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Comment == "" {
		t.Fatalf("expected non-empty comment draft, got %+v", resp)
	}
}

func TestRejectReviewMissingComment(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		rejectFn: func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
			return nil, service.ErrCommentRequired
		},
	})

	body := []byte(`{"comment":"   "}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews/1/reject", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestRejectReviewSuccess(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		rejectFn: func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
			return &model.PRActionResult{
				ReviewID: id,
				Action:   "rejected",
				Message:  "Requested changes and posted comment on GitHub.",
			}, nil
		},
	})

	body := []byte(`{"comment":"please fix issues"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews/1/reject", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp model.PRActionResult
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Action != "rejected" {
		t.Fatalf("unexpected action: %+v", resp)
	}
}

func TestApproveReviewNotFound(t *testing.T) {
	router := setupTestRouter(&mockReviewService{
		approveFn: func(ctx context.Context, id int64, comment string) (*model.PRActionResult, error) {
			return nil, service.ErrReviewNotFound
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews/99/approve", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusNotFound, rec.Body.String())
	}
}
