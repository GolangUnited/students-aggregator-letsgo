package main

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/config"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/parser/godev"
)

func GetParsers(pc []config.ParserConfig) (parsers []parser.ArticlesParser, err error) {
	for _, v := range pc {
		switch v.Name {
		case "go.dev":
			parsers = append(parsers, godev.NewParser(v.Url))
		default:
			return nil, fmt.Errorf("unsupported parser name %s", v.Name)
		}
	}
	return parsers, nil
}
