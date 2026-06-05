package md2html

import (
	"fmt"

	"md2html/internal/model"
	md2htmlRepo "md2html/internal/repository/tools/md2html"
)

type HistoryService interface {
	Save(userID int64, markdown, html, theme string) (int64, error)
	ListByUserID(userID int64, page, pageSize int) (*HistoryPageResult, error)
	DeleteByIDAndUserID(id, userID int64) error
}

type HistoryPageResult struct {
	Data  []model.Conversion
	Total int64
	Page  int
	Size  int
}

type historyService struct {
	historyRepo md2htmlRepo.HistoryRepository
}

func NewHistoryService(historyRepo md2htmlRepo.HistoryRepository) HistoryService {
	return &historyService{historyRepo: historyRepo}
}

func (s *historyService) Save(userID int64, markdown, html, theme string) (int64, error) {
	id, err := s.historyRepo.Save(userID, markdown, html, theme)
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

	return &HistoryPageResult{
		Data:  items,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

func (s *historyService) DeleteByIDAndUserID(id, userID int64) error {
	if err := s.historyRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return fmt.Errorf("delete history: %w", err)
	}
	return nil
}
