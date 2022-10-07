package main

import (
	"github.com/indikator/aggregator_lets_go/cmd/webservice"
	"github.com/indikator/aggregator_lets_go/model"
	"time"
)

var news = []model.Article{
	{
		Title:       "1",
		Created:     time.Now(),
		Description: "Short description of the article1",
		URL:         "http://funlink1.com",
		Author:      "m.mikhailov",
	},
	{
		Title:       "2",
		Created:     time.Now(),
		Description: "Short description of the article2",
		URL:         "http://funlink2.com",
		Author:      "m.mikhailov",
	},
	{
		Title:       "3",
		Created:     time.Now(),
		Description: "Short description of the article3",
		URL:         "http://funlink3.com",
		Author:      "m.mikhailov",
	},
}

func main() {
	err := webservice.RunServer(news)
	if err != nil {
		panic(err)
	}
}
