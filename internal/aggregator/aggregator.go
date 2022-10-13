package aggregator

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Execute() {
	c := config.NewConfig()

	err := c.Read("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	mongo.DBinit(c.Database.Url)

	parsers, err := GetParsers(c.Parsers)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range parsers {
		articles, err := v.ParseAll()

		if err != nil {
			log.Println(err)
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
	fmt.Println(c)
}
