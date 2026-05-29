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

const migrationFile = "migrations/001_create_reviews.sql"

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
	data, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	statements := strings.Split(string(data), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("execute migration: %w", err)
		}
	}
	return nil
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Save(ctx context.Context, result *model.AIReviewResult) (int64, error) {
	issues, err := json.Marshal(result.Issues)
	if err != nil {
		return 0, fmt.Errorf("marshal issues: %w", err)
	}

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO reviews (pr_title, pr_number, issues) VALUES (?, ?, ?)`,
		result.PRTitle, result.PRNumber, issues,
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

func (r *ReviewRepository) GetByID(ctx context.Context, id int64) (*model.ReviewRecord, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, pr_title, pr_number, issues, created_at FROM reviews WHERE id = ?`,
		id,
	)

	var record model.ReviewRecord
	var issuesJSON []byte
	var createdAt time.Time
	if err := row.Scan(&record.ID, &record.PRTitle, &record.PRNumber, &issuesJSON, &createdAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan review: %w", err)
	}

	if err := json.Unmarshal(issuesJSON, &record.Issues); err != nil {
		return nil, fmt.Errorf("unmarshal issues: %w", err)
	}
	record.CreatedAt = createdAt.UTC().Format(time.RFC3339)
	return &record, nil
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
		countQuery = `SELECT COUNT(*) FROM reviews WHERE pr_number = ?`
		countArgs = []any{prNumber}
	} else {
		countQuery = `SELECT COUNT(*) FROM reviews`
	}
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count reviews: %w", err)
	}

	var rows *sql.Rows
	var err error
	if prNumber > 0 {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, pr_title, pr_number, issues, created_at FROM reviews WHERE pr_number = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
			prNumber, limit, offset,
		)
	} else {
		rows, err = r.db.QueryContext(ctx,
			`SELECT id, pr_title, pr_number, issues, created_at FROM reviews ORDER BY created_at DESC LIMIT ? OFFSET ?`,
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
		var issuesJSON []byte
		var createdAt time.Time
		if err := rows.Scan(&item.ID, &item.PRTitle, &item.PRNumber, &issuesJSON, &createdAt); err != nil {
			return nil, 0, fmt.Errorf("scan review list item: %w", err)
		}
		var issues []model.ReviewIssue
		if err := json.Unmarshal(issuesJSON, &issues); err != nil {
			return nil, 0, fmt.Errorf("unmarshal issues: %w", err)
		}
		item.IssueCount = len(issues)
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
