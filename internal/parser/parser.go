package parser

import (
	"time"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/model"
)

// articles parser interface
type ArticlesParser interface {
	ParseAfter(date time.Time) (articles []model.Article, err error)
}

type NewParserFunc func(cfg Config, l log.Log) ArticlesParser

var ParserDefinitions map[string]NewParserFunc = make(map[string]NewParserFunc, 0)

func RegisterParser(name string, newParserFunc NewParserFunc) {
	ParserDefinitions[name] = newParserFunc
}
