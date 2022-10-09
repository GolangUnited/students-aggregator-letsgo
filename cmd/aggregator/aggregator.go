package main

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/config"
)

func main() {
	c := config.NewConfig()

	err := c.Read("aggregator.yaml")

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
