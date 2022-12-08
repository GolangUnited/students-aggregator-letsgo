package db

import (
	"github.com/indikator/aggregator_lets_go/model"
)

// Db interface to write and read to/from db
type Db interface {
	Name() string
	Url() string
	WriteArticle(article *model.DBArticle) (*model.DBArticle, error)
	ReadArticles(nDays int) ([]model.DBArticle, error)
	InitDb() error
}

func ConvertToDbArticle(article model.Article) *model.DBArticle {
	return &model.DBArticle{
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

func ConvertToDbArticles(articles []model.Article) []model.DBArticle {
	var r []model.DBArticle

	for _, a := range articles {
		r = append(r, *(ConvertToDbArticle(a)))
	}

	return r
}

func ConvertFromDbArticles(articles []model.DBArticle) []model.Article {
	var r []model.Article

	for _, a := range articles {
		r = append(r, *(ConvertFromDbArticle(a)))
	}

	return r
}
