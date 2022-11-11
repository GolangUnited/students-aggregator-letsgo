package stub

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
)

type database struct {
	config   db.Config
	Articles []model.DBArticle
}

// NewDb create an instance of database
func NewDb(config db.Config) db.Db {
	return &database{config: config}
}

func (d *database) DBInit() error {
	if d.config.Name != "stub" {
		return fmt.Errorf("incorrect db name %s, expected %s", d.config.Name, "stub")
	}
	return nil
}

func (d *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	d.Articles = append(d.Articles, *article)

	return article, nil
}

func (d *database) ReadArticles(nDays int) ([]model.DBArticle, error) {
	return d.Articles, nil
}
