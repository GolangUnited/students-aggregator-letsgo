// Aggregator parses sites, collects last articles from sites, and saves articles in database
//
// Aggregator uses config.yaml file with application settings
package main

import (
	"log"

	"github.com/indikator/aggregator_lets_go/internal/aggregator"
	"github.com/indikator/aggregator_lets_go/internal/config"
)

func main() {
	c := config.NewConfig()

	c.SetDataFromFile("config.yaml")

	a := aggregator.NewAggregator()

	err := a.InitAllByConfig(c)

	if err != nil {
		log.Fatal(err)
	}

	err = a.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
