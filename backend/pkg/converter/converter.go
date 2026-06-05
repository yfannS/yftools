package converter

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type GoldmarkConverter struct {
	md goldmark.Markdown
}

func NewGoldmarkConverter() *GoldmarkConverter {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.NewTable(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	return &GoldmarkConverter{md: md}
}

func (c *GoldmarkConverter) Convert(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := c.md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}
