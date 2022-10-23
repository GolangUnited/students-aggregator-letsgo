package stub

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
)

type Db struct {
	config   db.Config
	Articles []model.DBArticle
}

func (db *Db) DBinit(uri string) {
	db.config.Url = uri
}

func (db *Db) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	db.Articles = append(db.Articles, *article)

	return article, nil
}

func (db *Db) ReadAllArticles() ([]model.DBArticle, error) {
	return db.Articles, nil
}
