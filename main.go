package main

import (
	"github.com/indikator/aggregator_lets_go/webservice"
	"time"
)

var news = []webservice.DB{
	{
		Id:          1,
		Datetime:    time.Now(),
		Description: "Short description of the article1",
		Link:        "http://funlink1.com",
		PageHtml:    "pages/page1.html",
	},
	{
		Id:          2,
		Datetime:    time.Now(),
		Description: "Short description of the article2",
		Link:        "http://funlink2.com",
		PageHtml:    "pages/page2.html",
	},
	{
		Id:          3,
		Datetime:    time.Now(),
		Description: "Short description of the article3",
		Link:        "http://funlink3.com",
		PageHtml:    "pages/page3.html",
	},
}

func main() {
	err := webservice.RunServer(news)
	if err != nil {
		panic(err)
	}
}
