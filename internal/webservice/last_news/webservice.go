package last_news

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
)

type webService struct {
	handle string
	port   uint16
}

func NewWebservice(config webservice.Config) webservice.Webservice {
	handle := config.Handle
	port := config.Port
	return &webService{
		handle: handle,
		port:   port,
	}
}

func (ws *webService) MessageHandler(db db.Db) http.Handler {
	news, err := db.ReadAllArticles()
	if err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		newsJson, _ := json.Marshal(news)
		_, err = w.Write(newsJson)
		if err != nil {
			return
		}
	})
}

func (ws *webService) RunServer(db db.Db) {
	mux := http.NewServeMux()
	mux.Handle(ws.handle, ws.MessageHandler(db))

	http.ListenAndServe(":"+strconv.Itoa(int(ws.port)), mux)
	log.Println("Listening...")
}
