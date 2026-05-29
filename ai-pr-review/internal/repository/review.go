package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ZGjie20/PR-check-r/ai-pr-review/config"
	"github.com/ZGjie20/PR-check-r/ai-pr-review/internal/model"
	_ "github.com/go-sql-driver/mysql"
)

var migrationFiles = []string{
	"migrations/001_create_reviews.sql",
	"migrations/002_create_pr_review_record.sql",
}

type ReviewRepository struct {
	db *sql.DB
}

func Connect(cfg *config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return db, nil
}

func InitSchema(ctx context.Context, db *sql.DB) error {
	for _, migrationFile := range migrationFiles {
		data, err := os.ReadFile(migrationFile)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", migrationFile, err)
		}

		statements := strings.Split(string(data), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return fmt.Errorf("execute migration %s: %w", migrationFile, err)
			}
		}
	}
	return nil
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Save(ctx context.Context, input *model.ReviewSaveInput) (int64, error) {
	result := input.Result
	reviewResult, err := json.Marshal(result.ReviewResult)
	if err != nil {
		return 0, fmt.Errorf("marshal review_result: %w", err)
	}

	reviewStatus := input.ReviewStatus
	if reviewStatus == "" {
		reviewStatus = "completed"
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO pr_review_record (
			pr_number, pr_title, repo_name, pr_url, ai_model, review_status,
			total_issues, high_issues, medium_issues, low_issues,
			review_result, raw_diff
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		result.PRNumber,
		result.PRTitle,
		nullIfEmpty(input.RepoName),
		input.PRURL,
		nullIfEmpty(input.AIModel),
		reviewStatus,
		result.TotalIssues,
		result.HighIssues,
		result.MediumIssues,
		result.LowIssues,
		reviewResult,
		result.RawDiff,
	)
	if err != nil {
		return 0, fmt.Errorf("insert review: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("last insert id: %w", err)
	}
	return id, nil
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func (r *ReviewRepository) GetByID(ctx context.Context, id int64) (*model.ReviewRecord, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, pr_title, pr_number, repo_name, pr_url, ai_model, review_status,
			total_issues, high_issues, medium_issues, low_issues,
			review_result, raw_diff, created_at
		FROM pr_review_record WHERE id = ?`,
		id,
	)

	record, err := scanReviewRecord(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan review: %w", err)
	}
	return record, nil
}

func (r *ReviewRepository) List(ctx context.Context, page, limit int, prNumber int) ([]model.ReviewListItem, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var total int
	var countQuery string
	var countArgs []any
	if prNumber > 0 {
		countQuery = `SELECT COUNT(*) FROM pr_review_record WHERE pr_number = ?`
		countArgs = []any{prNumber}
	} else {
		countQuery = `SELECT COUNT(*) FROM pr_review_record`
	}
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews: %w", err)
	}

	var rows *sql.Rows
	var err error
	if prNumber > 0 {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, pr_title, pr_number, total_issues, created_at
			FROM pr_review_record WHERE pr_number = ?
			ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			prNumber, limit, offset,
		)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, pr_title, pr_number, total_issues, created_at
			FROM pr_review_record
			ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			limit, offset,
		)
	}
	if err != nil {
		return nil, 0, fmt.Errorf("list reviews: %w", err)
	}
	defer rows.Close()

	var items []model.ReviewListItem
	for rows.Next() {
		var item model.ReviewListItem
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.PRTitle, &item.PRNumber, &item.TotalIssues, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("scan review list item: %w", err)
		}
		item.CreatedAt = createdAt.UTC().Format(time.RFC3339)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate reviews: %w", err)
	}
	if items == nil {
		items = []model.ReviewListItem{}
	}
	return items, total, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanReviewRecord(row rowScanner) (*model.ReviewRecord, error) {
	var record model.ReviewRecord
	var repoName, aiModel sql.NullString
	var reviewResultJSON []byte
	var createdAt time.Time

	if err := row.Scan(
		&record.ID,
		&record.PRTitle,
		&record.PRNumber,
		&repoName,
		&record.PRURL,
		&aiModel,
		&record.ReviewStatus,
		&record.TotalIssues,
		&record.HighIssues,
		&record.MediumIssues,
		&record.LowIssues,
		&reviewResultJSON,
		&record.RawDiff,
		&createdAt,
	); err != nil {
		return nil, err
	}

	if repoName.Valid {
		record.RepoName = repoName.String
	}
	if aiModel.Valid {
		record.AIModel = aiModel.String
	}
	if err := json.Unmarshal(reviewResultJSON, &record.ReviewResult); err != nil {
		return nil, fmt.Errorf("unmarshal review_result: %w", err)
	}
	record.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return &record, nil
}
