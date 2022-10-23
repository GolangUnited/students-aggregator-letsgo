package aggregator

import (
	"fmt"
	"log"

	aconfig "github.com/indikator/aggregator_lets_go/internal/aggregator/config"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aggregator struct {
	config  aconfig.Config
	parsers []parser.ArticlesParser
	db      db.Db
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) InitAllByConfig(config *config.Config) error {
	err := config.Read()

	if err != nil {
		return err
	}

	parsers, err := GetParsers(config.Parsers)

	if err != nil {
		return err
	}

	db, err := GetDb(config.Database)
	mongoDb := mongo.NewDb(a.config.Database.Url)

	if err != nil {
		return err
	}

	return a.Init(&config.Aggregator, parsers, db)
}

func (a *Aggregator) Init(config *aconfig.Config, parsers []parser.ArticlesParser, db db.Db) error {
	a.config = *config

	a.parsers = parsers

	a.db = db

	return nil
}

func (a *Aggregator) Execute() error {
	for _, v := range a.parsers {
		articles, err := v.ParseAll()

		if err != nil {
			return err
		}

		for _, a := range articles {
			id := primitive.NewObjectID()

			_, err = mongoDb.WriteArticle(&model.DBArticle{
				ID:          id,
				Title:       a.Title,
				Created:     a.Created,
				Author:      a.Author,
				Description: a.Description,
				URL:         a.URL,
			})

			if err != nil {
				log.Println(err)
			}
		}

		fmt.Println(articles)
	}

	fmt.Println(a.config)

	return nil
}
