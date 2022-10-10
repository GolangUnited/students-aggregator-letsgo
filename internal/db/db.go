package db

import (
	"github.com/indikator/aggregator_lets_go/model"
)

// Db interface to write and read to/from db
type Db interface {
	DBinit(uri string)
	WriteArticle(article *model.DBArticle) (*model.DBArticle, error)
	ReadAllArticles() ([]model.DBArticle, error)
}
