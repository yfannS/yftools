package md2html

import (
	"md2html/pkg/converter"
)

type ConvertService interface {
	Convert(markdown string) (rawHTML, formattedHTML string, err error)
}

type convertService struct {
	converter *converter.GoldmarkConverter
}

func NewConvertService(c *converter.GoldmarkConverter) ConvertService {
	return &convertService{converter: c}
}

func (s *convertService) Convert(markdown string) (string, string, error) {
	rawHTML, err := s.converter.Convert(markdown)
	if err != nil {
		return "", "", err
	}

	// formattedHTML = rawHTML for now (Goldmark output is already formatted)
	formatted := rawHTML

	return rawHTML, formatted, nil
}
