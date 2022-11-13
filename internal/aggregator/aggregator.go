package aggregator

import (
	"fmt"
	"time"

	aconfig "github.com/indikator/aggregator_lets_go/internal/aggregator/config"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aggregator struct {
	config            aconfig.Config
	parsers           []parser.ArticlesParser
	db                db.Db
	lastCheckDatetime time.Time
}

const (
	daysAgo = 1
)

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
	a.lastCheckDatetime = time.Now().AddDate(0, -daysAgo, 0)

	var errAll error = nil

	for _, v := range a.parsers {
		articles, err := v.ParseAfter(a.lastCheckDatetime)

		if err != nil {
			if errAll != nil {
				errAll = fmt.Errorf("%v; %v", errAll, err)
			} else {
				errAll = err
			}
			continue
		}

		for _, article := range articles {
			id := primitive.NewObjectID()

			_, err = a.db.WriteArticle(&model.DBArticle{
				ID:          id,
				Title:       article.Title,
				Created:     article.Created,
				Author:      article.Author,
				Description: article.Description,
				URL:         article.URL,
			})

			if err != nil {
				if errAll != nil {
					errAll = fmt.Errorf("%v; %v", errAll, err)
				} else {
					errAll = err
				}
			}
		}
	}

	return errAll
}
