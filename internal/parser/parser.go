package parser

import (
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

/// parser interface
type NewsParser interface {
	ParseAll() (posts []model.Post, err error)
	ParseBefore(date time.Time) (posts []model.Post, err error)
	ParseBeforeN(date time.Time, n int) (posts []model.Post, err error)
}
