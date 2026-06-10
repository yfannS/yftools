package jsonformatter

import (
	"database/sql"
	"fmt"

	"md2html/internal/model"
)

type JsonHistoryRepository interface {
	Save(userID int64, title, input, output string, charCount int) (int64, error)
	ListByUserID(userID int64, page, pageSize int) ([]model.JsonHistoryListItem, int64, error)
	FindByIDAndUserID(id, userID int64) (*model.JsonHistoryDetail, error)
	DeleteByIDAndUserID(id, userID int64) error
	UpdateTitle(id, userID int64, title string) error
}

type jsonHistoryRepository struct {
	db *sql.DB
}

func NewJsonHistoryRepository(db *sql.DB) JsonHistoryRepository {
	return &jsonHistoryRepository{db: db}
}

func (r *jsonHistoryRepository) Save(userID int64, title, input, output string, charCount int) (int64, error) {
	result, err := r.db.Exec(
		"INSERT INTO json_format_history (user_id, title, input, output, char_count, is_delete, created_at, updated_at) VALUES (?, ?, ?, ?, ?, 0, NOW(), NOW())",
		userID, title, input, output, charCount,
	)
	if err != nil {
		return 0, fmt.Errorf("save json history: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	return id, nil
}

func (r *jsonHistoryRepository) ListByUserID(userID int64, page, pageSize int) ([]model.JsonHistoryListItem, int64, error) {
	var total int64
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM json_format_history WHERE user_id = ? AND is_delete = 0",
		userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count json history: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		"SELECT id, title, char_count, created_at, updated_at FROM json_format_history WHERE user_id = ? AND is_delete = 0 ORDER BY updated_at DESC LIMIT ? OFFSET ?",
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list json history: %w", err)
	}
	defer rows.Close()

	var items []model.JsonHistoryListItem
	for rows.Next() {
		var item model.JsonHistoryListItem
		if err := rows.Scan(&item.ID, &item.Title, &item.CharCount, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan json history: %w", err)
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *jsonHistoryRepository) FindByIDAndUserID(id, userID int64) (*model.JsonHistoryDetail, error) {
	detail := &model.JsonHistoryDetail{}

	err := r.db.QueryRow(
		"SELECT id, title, input, output, char_count, created_at, updated_at FROM json_format_history WHERE id = ? AND user_id = ? AND is_delete = 0",
		id, userID,
	).Scan(&detail.ID, &detail.Title, &detail.Input, &detail.Output, &detail.CharCount, &detail.CreatedAt, &detail.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find json history by id: %w", err)
	}

	return detail, nil
}

func (r *jsonHistoryRepository) DeleteByIDAndUserID(id, userID int64) error {
	result, err := r.db.Exec(
		"UPDATE json_format_history SET is_delete = 1, updated_at = NOW() WHERE id = ? AND user_id = ? AND is_delete = 0",
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("soft delete json history: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *jsonHistoryRepository) UpdateTitle(id, userID int64, title string) error {
	result, err := r.db.Exec(
		"UPDATE json_format_history SET title = ?, updated_at = NOW() WHERE id = ? AND user_id = ? AND is_delete = 0",
		title, id, userID,
	)
	if err != nil {
		return fmt.Errorf("update json history title: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
