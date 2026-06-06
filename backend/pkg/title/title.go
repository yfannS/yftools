package title

import (
	"strings"
	"unicode/utf8"
)

const maxTitleLen = 200

// ExtractTitle 从 Markdown 文本中提取文档标题
// 优先级：# > ## > ### > #### > ##### > ###### > 首行非空 > "未命名文档"
func ExtractTitle(markdown string) string {
	lines := strings.Split(markdown, "\n")

	// 优先级 1-6：匹配 # ~ ###### 标题
	for level := 1; level <= 6; level++ {
		prefix := strings.Repeat("#", level) + " "
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, prefix) {
				titleText := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
				if titleText != "" {
					return truncate(titleText, maxTitleLen)
				}
			}
		}
	}

	// 优先级 7：第一个非空行（最多取 60 字符）
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			// 跳过纯标记行（如 ---, ***, ```）
			if isDecorationLine(trimmed) {
				continue
			}
			if utf8.RuneCountInString(trimmed) > 60 {
				return truncate(trimmed, 60) + "..."
			}
			return trimmed
		}
	}

	return "未命名文档"
}

// truncate 截断字符串到 maxRunes 个 rune
func truncate(s string, maxRunes int) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxRunes])
}

// isDecorationLine 判断是否为纯装饰性行（不应作为标题）
func isDecorationLine(line string) bool {
	// 水平分隔线：---, ***, ___
	trimmed := line
	allSame := true
	ch := trimmed[0]
	for _, c := range trimmed {
		if c != rune(ch) {
			allSame = false
			break
		}
	}
	if allSame && len(trimmed) >= 3 {
		return true
	}
	// 代码块开始
	if strings.HasPrefix(trimmed, "```") {
		return true
	}
	return false
}
