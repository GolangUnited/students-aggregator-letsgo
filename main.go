package main

import (
	"github.com/indikator/aggregator_lets_go/webservice"
	"log"
	"net/http"
	"time"
)

var news = []webservice.DB{
	webservice.DB{
		Id:          1,
		Datetime:    time.Now(),
		Description: "Short description of the article1",
		Link:        "http://funlink1.com",
		PageHtml:    "pages/page1.html",
	},
	webservice.DB{
		Id:          2,
		Datetime:    time.Now(),
		Description: "Short description of the article2",
		Link:        "http://funlink2.com",
		PageHtml:    "pages/page2.html",
	},
	webservice.DB{
		Id:          3,
		Datetime:    time.Now(),
		Description: "Short description of the article3",
		Link:        "http://funlink3.com",
		PageHtml:    "pages/page3.html",
	},
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/last_news", webservice.MessageHandler(news))
	log.Println("Listening...")
	http.ListenAndServe(":8080", mux)
}
