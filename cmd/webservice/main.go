package main

import (
	"log"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
)

const (
	configFilePath = "./configs/config.yaml"
)

func main() {
	cfg := config.NewConfig()
	if err := cfg.SetDataFromFile(configFilePath); err != nil {
		log.Fatal(err)
	}

	ws := last_news.NewWebservice()

	if err := ws.InitAllByConfig(cfg); err != nil {
		log.Fatal(err)
	}

	ws.RunServer()
}
