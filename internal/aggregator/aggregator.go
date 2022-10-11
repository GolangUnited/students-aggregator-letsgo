package aggregator

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/internal/config"
)

func Execute() {
	c := config.NewConfig()

	err := c.Read("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	parsers, err := GetParsers(c.Parsers)

	if err != nil {
		log.Fatal(err)
	}

	for _, v := range parsers {
		articles, err := v.ParseAll()

		if err != nil {
			log.Println(err)
		}

		fmt.Println(articles)
	}
	fmt.Println(c)
}
