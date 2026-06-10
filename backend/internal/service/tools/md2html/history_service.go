package md2html

import (
	"fmt"
	"unicode/utf8"

	"md2html/internal/model"
	md2htmlRepo "md2html/internal/repository/tools/md2html"
	"md2html/pkg/title"
)

type HistoryService interface {
	Save(userID int64, markdown, html, theme string) (int64, error)
	ListByUserID(userID int64, page, pageSize int) (*HistoryPageResult, error)
	GetDetail(id, userID int64) (*model.HistoryDetail, error)
	DeleteByIDAndUserID(id, userID int64) error
	RenameTitle(id, userID int64, title string) error
}

type HistoryPageResult struct {
	Data       []model.HistoryListItem
	Total      int64
	Page       int
	Size       int
	TotalPages int
}

type historyService struct {
	historyRepo md2htmlRepo.HistoryRepository
}

func NewHistoryService(historyRepo md2htmlRepo.HistoryRepository) HistoryService {
	return &historyService{historyRepo: historyRepo}
}

func (s *historyService) Save(userID int64, markdown, html, theme string) (int64, error) {
	// 自动提取标题
	docTitle := title.ExtractTitle(markdown)
	// 计算 Unicode 字符数
	charCount := utf8.RuneCountInString(markdown)

	id, err := s.historyRepo.Save(userID, docTitle, markdown, html, charCount, theme)
	if err != nil {
		return 0, fmt.Errorf("save history: %w", err)
	}
	return id, nil
}

func (s *historyService) ListByUserID(userID int64, page, pageSize int) (*HistoryPageResult, error) {
	items, total, err := s.historyRepo.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list history: %w", err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &HistoryPageResult{
		Data:       items,
		Total:      total,
		Page:       page,
		Size:       pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *historyService) GetDetail(id, userID int64) (*model.HistoryDetail, error) {
	detail, err := s.historyRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("get detail: %w", err)
	}
	if detail == nil {
		return nil, fmt.Errorf("记录不存在")
	}
	return detail, nil
}

func (s *historyService) DeleteByIDAndUserID(id, userID int64) error {
	if err := s.historyRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return fmt.Errorf("delete history: %w", err)
	}
	return nil
}

func (s *historyService) RenameTitle(id, userID int64, title string) error {
	if title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if err := s.historyRepo.UpdateTitle(id, userID, title); err != nil {
		return fmt.Errorf("rename history: %w", err)
	}
	return nil
}
