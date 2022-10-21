package aggregator

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aggregator struct {
	config config.Config
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Init(config *config.Config) error {
	a.config = *config

	return nil
}

func (a *Aggregator) Execute() error {
	err := a.config.Read()

	if err != nil {
		return err
	}

	parsers, err := GetParsers(a.config.Parsers)

	mongo.DBinit(a.config.Database.Url)

	if err != nil {
		return err
	}

	for _, v := range parsers {
		articles, err := v.ParseAll()

		if err != nil {
			return err
		}

		for _, a := range articles {
			id := primitive.NewObjectID()

			_, err = mongo.WriteArticle(&model.DBArticle{
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
