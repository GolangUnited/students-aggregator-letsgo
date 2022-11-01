package last_news

import (
	"encoding/json"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"github.com/indikator/aggregator_lets_go/model"
	"log"
	"net/http"
	"strconv"
	"time"
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

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.RFC3339)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, time.RFC3339)
	b = append(b, '"')
	return b, nil
}

func (ws *webService) MessageHandler(db db.Db) http.Handler {

	news, err := db.ReadAllArticles()
	if err != nil {
		log.Fatal(err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		for i, value := range news {
			news[i] = model.DBArticle{
				ID:          value.ID,
				Title:       value.Title,
				Created:     value.Created,
				Author:      value.Author,
				Description: value.Description,
				URL:         value.URL}
		}
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
