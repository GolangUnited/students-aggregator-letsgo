package parser

import (
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

/// articles parser interface
type ArticlesParser interface {
	ParseAll() (posts []model.Article, err error)
	ParseAfter(date time.Time) (posts []model.Article, err error)
	ParseAfterN(date time.Time, n int) (posts []model.Article, err error)
}

type NewParserFunc func(string) ArticlesParser

var ParserDefinitions map[string]NewParserFunc = make(map[string]NewParserFunc, 0)

func RegisterParser(name string, newParserFunc NewParserFunc) {
	ParserDefinitions[name] = newParserFunc
}
