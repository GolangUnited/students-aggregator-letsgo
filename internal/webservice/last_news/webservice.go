package last_news

import (
	"encoding/json"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"log"
	"net/http"
	"strconv"
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

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		news, err := db.ReadArticles(7) // hardcode 1 week
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		newsJson, err := json.Marshal(news)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(newsJson)
	})
}

func (ws *webService) RunServer(db db.Db) {
	mux := http.NewServeMux()
	mux.Handle(ws.handle, ws.MessageHandler(db))
	log.Println("Listening...")
	http.ListenAndServe(":"+strconv.Itoa(int(ws.port)), mux)
}
