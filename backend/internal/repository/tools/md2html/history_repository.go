package md2html

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"md2html/internal/model"
)

type HistoryRepository interface {
	// Save 保存记录（INSERT 或 UPSERT），返回记录 ID
	Save(userID int64, title, markdown, html string, charCount int, theme string) (int64, error)
	// ListByUserID 列表查询（轻量，不含 markdown/html）
	ListByUserID(userID int64, page, pageSize int) ([]model.HistoryListItem, int64, error)
	// FindByIDAndUserID 详情查询（含完整 markdown）
	FindByIDAndUserID(id, userID int64) (*model.HistoryDetail, error)
	// DeleteByIDAndUserID 删除
	DeleteByIDAndUserID(id, userID int64) error
	// FindRecentByUserIDAndTitle 查找用户最近一条相同标题的记录（UPSERT 判断）
	FindRecentByUserIDAndTitle(userID int64, title string) (*model.Conversion, error)
}

type historyRepository struct {
	db *sql.DB
}

func NewHistoryRepository(db *sql.DB) HistoryRepository {
	return &historyRepository{db: db}
}

func (r *historyRepository) Save(userID int64, title, markdown, html string, charCount int, theme string) (int64, error) {
	if theme == "" {
		theme = "default"
	}

	// UPSERT 策略：同用户 + 同标题 + 5分钟内 → 更新
	existing, err := r.FindRecentByUserIDAndTitle(userID, title)
	if err == nil && existing != nil {
		// 5分钟窗口内，更新现有记录
		log.Printf("[Repo.Save] UPSERT update id=%d, markdown_len=%d", existing.ID, len(markdown))
		_, err = r.db.Exec(
			"UPDATE md2html_conversions SET markdown = ?, html = ?, char_count = ?, theme = ?, updated_at = NOW() WHERE id = ?",
			markdown, sql.NullString{String: html, Valid: html != ""}, charCount, theme, existing.ID,
		)
		if err != nil {
			return 0, fmt.Errorf("update conversion: %w", err)
		}
		return existing.ID, nil
	}

	// 插入新记录
	var htmlVal interface{}
	if html != "" {
		htmlVal = html
	} else {
		htmlVal = nil
	}

	result, err := r.db.Exec(
		"INSERT INTO md2html_conversions (user_id, title, markdown, html, char_count, theme, is_delete, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, 0, NOW(), NOW())",
		userID, title, markdown, htmlVal, charCount, theme,
	)
	if err != nil {
		log.Printf("[Repo.Save] INSERT failed: %v", err)
		return 0, fmt.Errorf("save conversion: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}

	log.Printf("[Repo.Save] INSERT success id=%d, markdown_len=%d", id, len(markdown))
	return id, nil
}

func (r *historyRepository) ListByUserID(userID int64, page, pageSize int) ([]model.HistoryListItem, int64, error) {
	var total int64
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM md2html_conversions WHERE user_id = ? AND is_delete = 0",
		userID,
	).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count conversions: %w", err)
	}

	offset := (page - 1) * pageSize
	rows, err := r.db.Query(
		"SELECT id, title, char_count, theme, created_at, updated_at FROM md2html_conversions WHERE user_id = ? AND is_delete = 0 ORDER BY updated_at DESC LIMIT ? OFFSET ?",
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list conversions: %w", err)
	}
	defer rows.Close()

	var items []model.HistoryListItem
	for rows.Next() {
		var item model.HistoryListItem
		if err := rows.Scan(&item.ID, &item.Title, &item.CharCount, &item.Theme, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan conversion: %w", err)
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *historyRepository) FindByIDAndUserID(id, userID int64) (*model.HistoryDetail, error) {
	detail := &model.HistoryDetail{}
	var htmlVal sql.NullString

	err := r.db.QueryRow(
		"SELECT id, title, markdown, html, char_count, theme, created_at, updated_at FROM md2html_conversions WHERE id = ? AND user_id = ? AND is_delete = 0",
		id, userID,
	).Scan(&detail.ID, &detail.Title, &detail.Markdown, &htmlVal, &detail.CharCount, &detail.Theme, &detail.CreatedAt, &detail.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find by id: %w", err)
	}

	if htmlVal.Valid {
		detail.HTML = htmlVal.String
	}

	return detail, nil
}

func (r *historyRepository) DeleteByIDAndUserID(id, userID int64) error {
	result, err := r.db.Exec(
		"UPDATE md2html_conversions SET is_delete = 1, updated_at = NOW() WHERE id = ? AND user_id = ? AND is_delete = 0",
		id, userID,
	)
	if err != nil {
		return fmt.Errorf("soft delete conversion: %w", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *historyRepository) FindRecentByUserIDAndTitle(userID int64, title string) (*model.Conversion, error) {
	c := &model.Conversion{}
	var updatedAt sql.NullTime

	err := r.db.QueryRow(
		"SELECT id, user_id, title, updated_at FROM md2html_conversions WHERE user_id = ? AND title = ? AND is_delete = 0 AND updated_at > ? ORDER BY updated_at DESC LIMIT 1",
		userID, title, time.Now().Add(-5*time.Minute),
	).Scan(&c.ID, &c.UserID, &c.Title, &updatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("find recent by title: %w", err)
	}

	if updatedAt.Valid {
		c.UpdatedAt = updatedAt.Time
	}

	return c, nil
}
