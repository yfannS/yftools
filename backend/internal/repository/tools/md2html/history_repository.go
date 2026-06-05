package md2html

import (
	"database/sql"
	"fmt"

	"md2html/internal/model"
)

type HistoryRepository interface {
	Save(userID int64, markdown, html, theme string) (int64, error)
	ListByUserID(userID int64, page, pageSize int) ([]model.Conversion, int64, error)
	DeleteByIDAndUserID(id, userID int64) error
}

type historyRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Save(userID int64, markdown, html, theme string) (int64, error) {
	if theme == "" {
		theme = "default"
	}

	result, err := r.db.Exec(
		"INSERT INTO md2html_conversions (user_id, markdown, html, theme, created_at) VALUES (?, ?, ?, ?, NOW())",
		userID, markdown, html, theme,
	)
	if err != nil {
		return 0, fmt.Errorf("save conversion: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (r *historyRepository) ListByUserID(userID int64, page, pageSize int) ([]model.Conversion, int64, error) {
	var total int64
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM md2html_conversions WHERE user_id = ?",
		userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count conversions: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		"SELECT id, user_id, markdown, html, theme, created_at FROM md2html_conversions WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?",
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list conversions: %w", err)
	}
	defer rows.Close()

	var items []model.Conversion
	for rows.Next() {
		var c model.Conversion
		if err := rows.Scan(&c.ID, &c.UserID, &c.Markdown, &c.HTML, &c.Theme, &c.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan conversion: %w", err)
		}
		items = append(items, c)
	}

	return items, total, nil
}

func (r *historyRepository) DeleteByIDAndUserID(id, userID int64) error {
	result, err := r.db.Exec(
		"DELETE FROM md2html_conversions WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("delete conversion: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
