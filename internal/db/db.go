package db

import (
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/mongo"
)

// Db interface to write and read to/from db
type Db interface {
	DBinit(uri string) *mongo.Collection
	WriteArticle(article *model.Article, collection *mongo.Collection) error
	ReadAllArticles(collection *mongo.Collection) ([]*model.Article, error)
}
