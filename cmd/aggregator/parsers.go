package main

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/parser"

	_ "github.com/indikator/aggregator_lets_go/internal/parser/autoregister"
)

func GetParsers(pc []parser.Config) (parsers []parser.ArticlesParser, err error) {
	for _, v := range pc {
		newParserFunc := parser.ParserDefinitions[v.Name]
		if newParserFunc == nil {
			return nil, fmt.Errorf("unsupported parser name %s", v.Name)
		}

		parsers = append(parsers, newParserFunc(v.Url))
	}
	return parsers, nil
}
