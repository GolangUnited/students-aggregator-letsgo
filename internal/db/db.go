package db

import (
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Db interface to write and read to/from db
type Db interface {
	Name() string
	Url() string
	WriteArticle(article *model.DBArticle) (*model.DBArticle, error)
	ReadArticles(nDays int) ([]model.DBArticle, error)
	InitDb() error
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

func ConvertToDbArticle(article model.Article) *model.DBArticle {
	return &model.DBArticle{
		ID:          primitive.NewObjectID(),
		Title:       article.Title,
		Created:     article.Created,
		Author:      article.Author,
		Description: article.Description,
		URL:         article.URL,
	}
}

func ConvertFromDbArticle(article model.DBArticle) *model.Article {
	return &model.Article{
		Title:       article.Title,
		Created:     article.Created,
		Author:      article.Author,
		Description: article.Description,
		URL:         article.URL,
	}
}
