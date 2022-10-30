package main

import (
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
	"log"
)

func main() {

	c := config.NewConfig()
	err := c.SetDataFromFile("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = c.Read()
	if err != nil {
		log.Fatal(err)
	}

	db_ := mongo.NewDb(c.Database)
	err = db_.DBInit()
	if err != nil {
		log.Fatal(err)
	}

	ws := last_news.NewWebservice(c.WebService)

	ws.RunServer(db_)
}
