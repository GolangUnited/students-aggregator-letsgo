package parser

import (
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

/// articles parser interface
type ArticlesParser interface {
	ParseAll() (articles []model.Article, err error)
	ParseAfter(date time.Time) (articles []model.Article, err error)
	ParseAfterN(date time.Time, n int) (articles []model.Article, err error)
}
