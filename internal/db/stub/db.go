package stub

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	WriteArticleErrorUrl = "stub://localhost:error/"
	InitDbErrorUrl       = ""
)

type database struct {
	config   db.Config
	log      log.Log
	nextId   int
	Articles []model.DBArticle
}

// NewDb create an instance of database
func NewDb(config db.Config, l log.Log) db.Db {
	return &database{config: config, log: l, nextId: 1}
}

func (d *database) Name() string {
	return d.config.Name
}

func (d *database) Url() string {
	return d.config.Url
}

func (d *database) InitDb() error {
	if d.config.Url == InitDbErrorUrl {
		return fmt.Errorf("incorrect db url")
	}

	return nil
}

func (d *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	if d.config.Url == WriteArticleErrorUrl {
		return nil, fmt.Errorf("incorrect db url")
	}

	if article.ID == nil {
		article.ID = d.nextId
		d.nextId++
	}
	d.Articles = append(d.Articles, *article)

	return article, nil
}

func (d *database) ReadArticles(nDays int) ([]model.DBArticle, error) {
	if nDays < 1 {
		return nil, fmt.Errorf("invalid number of days %d", nDays)
	}
	return d.Articles, nil
}
