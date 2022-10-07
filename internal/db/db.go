package db

import (
	"github.com/indikator/aggregator_lets_go/model"
)

// DbWriter interface to write and read to/from db
type DbWriter interface {
	DBinit(uri string)
	WriteArticle(article *model.Article) error
	ReadAll() ([]*model.Article, error)
}
