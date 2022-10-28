package main

import (
	"log"

	"github.com/indikator/aggregator_lets_go/internal/aggregator"
	"github.com/indikator/aggregator_lets_go/internal/config"
)

func main() {

	c := config.NewConfig()

	err := c.SetDataFromFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	a := aggregator.NewAggregator()

	err = a.InitAllByConfig(c)

	if err != nil {
		log.Fatal(err)
	}

	err = a.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
