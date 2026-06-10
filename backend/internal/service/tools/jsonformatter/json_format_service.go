package jsonformatter

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf8"

	"md2html/internal/model"
	jsonRepo "md2html/internal/repository/tools/jsonformatter"
)

type JsonFormatService interface {
	Format(input string, indent int, minify bool) (*model.JsonFormatResponse, error)
	Validate(input string) (*model.JsonValidateResponse, error)
	SaveHistory(userID int64, input, output string) (int64, error)
	ListHistory(userID int64, page, pageSize int) (*JsonHistoryPageResult, error)
	GetHistoryDetail(id, userID int64) (*model.JsonHistoryDetail, error)
	DeleteHistory(id, userID int64) error
	RenameHistory(id, userID int64, title string) error
}

type JsonHistoryPageResult struct {
	Data       []model.JsonHistoryListItem
	Total      int64
	Page       int
	Size       int
	TotalPages int
}

type jsonFormatService struct {
	jsonRepo jsonRepo.JsonHistoryRepository
}

func NewJsonFormatService(jsonRepo jsonRepo.JsonHistoryRepository) JsonFormatService {
	return &jsonFormatService{jsonRepo: jsonRepo}
}

func (s *jsonFormatService) Format(input string, indent int, minify bool) (*model.JsonFormatResponse, error) {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return nil, fmt.Errorf("无效的 JSON: %w", err)
	}

	var output []byte
	var err error

	if minify {
		output, err = json.Marshal(obj)
	} else {
		if indent <= 0 {
			indent = 2
		}
		prefix := ""
		indentStr := strings.Repeat(" ", indent)
		output, err = json.MarshalIndent(obj, prefix, indentStr)
	}

	if err != nil {
		return nil, fmt.Errorf("格式化失败: %w", err)
	}

	return &model.JsonFormatResponse{
		Output:   string(output),
		Size:     len(output),
		Minified: minify,
	}, nil
}

func (s *jsonFormatService) Validate(input string) (*model.JsonValidateResponse, error) {
	var obj interface{}
	decoder := json.NewDecoder(strings.NewReader(input))
	decoder.UseNumber()
	err := decoder.Decode(&obj)

	if err != nil {
		return &model.JsonValidateResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	keys, depth := countKeysAndDepth(obj)

	return &model.JsonValidateResponse{
		Valid: true,
		Keys:  keys,
		Depth: depth,
	}, nil
}

func (s *jsonFormatService) SaveHistory(userID int64, input, output string) (int64, error) {
	title := extractJsonTitle(input)
	charCount := utf8.RuneCountInString(input)

	id, err := s.jsonRepo.Save(userID, title, input, output, charCount)
	if err != nil {
		return 0, fmt.Errorf("save json history: %w", err)
	}
	return id, nil
}

func (s *jsonFormatService) ListHistory(userID int64, page, pageSize int) (*JsonHistoryPageResult, error) {
	items, total, err := s.jsonRepo.ListByUserID(userID, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("list json history: %w", err)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	return &JsonHistoryPageResult{
		Data:       items,
		Total:      total,
		Page:       page,
		Size:       pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *jsonFormatService) GetHistoryDetail(id, userID int64) (*model.JsonHistoryDetail, error) {
	detail, err := s.jsonRepo.FindByIDAndUserID(id, userID)
	if err != nil {
		return nil, fmt.Errorf("get json history detail: %w", err)
	}
	if detail == nil {
		return nil, fmt.Errorf("记录不存在")
	}
	return detail, nil
}

func (s *jsonFormatService) DeleteHistory(id, userID int64) error {
	if err := s.jsonRepo.DeleteByIDAndUserID(id, userID); err != nil {
		return fmt.Errorf("delete json history: %w", err)
	}
	return nil
}

func (s *jsonFormatService) RenameHistory(id, userID int64, title string) error {
	if title == "" {
		return fmt.Errorf("标题不能为空")
	}
	if err := s.jsonRepo.UpdateTitle(id, userID, title); err != nil {
		return fmt.Errorf("rename json history: %w", err)
	}
	return nil
}

// extractJsonTitle 从 JSON 中提取第一个有意义的字符串值作为标题
func extractJsonTitle(input string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err != nil {
		return "未命名 JSON"
	}

	return findFirstStringValue(obj, 0)
}

func findFirstStringValue(v interface{}, depth int) string {
	if depth > 3 {
		return "未命名 JSON"
	}

	switch val := v.(type) {
	case map[string]interface{}:
		// 优先查找常见标题字段
		for _, key := range []string{"name", "title", "label", "description", "id"} {
			if fv, ok := val[key]; ok {
				if s, ok := fv.(string); ok && s != "" {
					return s
				}
			}
		}
		// 取第一个字符串值
		for _, fv := range val {
			if s, ok := fv.(string); ok && s != "" {
				return s
			}
		}
		// 递归第一个子对象
		for _, fv := range val {
			if result := findFirstStringValue(fv, depth+1); result != "未命名 JSON" {
				return result
			}
		}
	case []interface{}:
		if len(val) > 0 {
			return findFirstStringValue(val[0], depth+1)
		}
	}

	return "未命名 JSON"
}

// countKeysAndDepth 统计 JSON 键数量和最大嵌套深度
func countKeysAndDepth(v interface{}) (keys, depth int) {
	return countKeys(v, 1)
}

func countKeys(v interface{}, d int) (keys, depth int) {
	switch val := v.(type) {
	case map[string]interface{}:
		keys = len(val)
		depth = d
		for _, fv := range val {
			ck, cd := countKeys(fv, d+1)
			keys += ck
			if cd > depth {
				depth = cd
			}
		}
	case []interface{}:
		depth = d
		for _, fv := range val {
			ck, cd := countKeys(fv, d+1)
			keys += ck
			if cd > depth {
				depth = cd
			}
		}
	default:
		depth = d
	}
	return
}
