package last_news

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
)

type webService struct {
	handle string
}

func NewWebservice(handle string) webservice.Webservice {
	return &webService{
		handle: handle,
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
		news, _ := json.Marshal(news)
		w.Write(news)
	})
}

func RunServer(ws webservice.Webservice, c config.Config, handle string) error {
	db := mongo.NewDb(c.Database.Url)
	err := db.DBInit()
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle(handle, ws.MessageHandler(db))
	err = http.ListenAndServe(":"+strconv.Itoa(int(c.WebService.Port)), mux)
	if err != nil {
		return fmt.Errorf("can't start the server: %w", err)
	}
	log.Println("Listening...")
	return nil
}
