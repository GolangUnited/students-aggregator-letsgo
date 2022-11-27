package stub

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/model"
)

type database struct {
	config   db.Config
	log      log.Log
	Articles []model.DBArticle
}

// NewDb create an instance of database
func NewDb(config db.Config, l log.Log) db.Db {
	return &database{config: config, log: l}
}

func (d *database) Name() string {
	return d.config.Name
}

func (d *database) Url() string {
	return d.config.Url
}

func (d *database) DBInit() error {
	return nil
}

func (d *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	d.Articles = append(d.Articles, *article)

	return article, nil
}

func (d *database) ReadArticles(nDays int) ([]model.DBArticle, error) {
	return d.Articles, nil
}
