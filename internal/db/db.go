package db

import (
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Db interface to write and read to/from db
type Db interface {
	WriteArticle(article *model.DBArticle) (*model.DBArticle, error)
	ReadAllArticles() ([]model.DBArticle, error)
	DBInit()
}

func ConvertArticle(article model.Article) *model.DBArticle {
	return &model.DBArticle{
		ID:          primitive.NewObjectID(),
		Title:       article.Title,
		Created:     article.Created,
		Author:      article.Author,
		Description: article.Description,
		URL:         article.URL,
	}
}
