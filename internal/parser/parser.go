package parser

import (
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

/// articles parser interface
type ArticlesParser interface {
	ParseAll() (posts []model.Article, err error)
	ParseBefore(date time.Time) (posts []model.Article, err error)
	ParseBeforeN(date time.Time, n int) (posts []model.Article, err error)
}
