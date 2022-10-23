package stub

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
)

type database struct {
	config   db.Config
	Articles []model.DBArticle
}

// NewDb create an instance of database
func NewDb(URL string) db.Db {
	return &database{}
}

func (d *database) DBInit(uri string) {
	d.config.Url = uri
}

func (d *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	d.Articles = append(d.Articles, *article)

	return article, nil
}

func (d *database) ReadAllArticles() ([]model.DBArticle, error) {
	return d.Articles, nil
}
