package aggregator

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/parser"

	_ "github.com/indikator/aggregator_lets_go/internal/parser/autoregister"
)

const (
	unsupportedParserNameErrorTemplate = "unsupported parser name %s. Read requirements from https://github.com/indikator/aggregator_lets_go/blob/main/README.md"
)

type UnsupportedParserNameError struct {
	Text string
}

func (e *UnsupportedParserNameError) Error() string {
	return e.Text
}

func GetParsers(pc []parser.Config, l log.Log) (parsers []parser.ArticlesParser, err error) {
	for _, v := range pc {
		newParserFunc := parser.ParserDefinitions[v.Name]
		if newParserFunc == nil {
			return nil, &UnsupportedParserNameError{
				Text: fmt.Sprintf(unsupportedParserNameErrorTemplate, v.Name),
			}
		}

		cfg := parser.Config{URL: v.URL}
		parsers = append(parsers, newParserFunc(cfg, l))
	}
	return parsers, nil
}
