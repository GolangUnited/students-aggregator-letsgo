// Aggregator parses sites, collects last articles from sites, and saves articles in database
//
// Aggregator uses config.yaml file with application settings
package main

import (
	"log"

	"github.com/indikator/aggregator_lets_go/internal/aggregator"
	"github.com/indikator/aggregator_lets_go/internal/config"
)

const (
	configFilePath = "./configs/config.yaml"
	//configFilePath = "../../configs/config.yaml"
)

func main() {
	c := config.NewConfig()

	err := c.SetDataFromFile(configFilePath)
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
