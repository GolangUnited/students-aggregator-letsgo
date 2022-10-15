package last_news

import (
	"encoding/json"
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"log"
	"net/http"
	"strconv"
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

func RunServer(ws webservice.Webservice, handle string) error {
	c := config.NewConfig()
	err := c.Read("etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	db := mongo.NewDb(c.Database.Url)
	mux := http.NewServeMux()
	mux.Handle(handle, ws.MessageHandler(db))
	err = http.ListenAndServe(":"+strconv.Itoa(int(c.WebService.Port)), mux)
	if err != nil {
		return fmt.Errorf("Can't start the server")
	}
	log.Println("Listening...")
	return nil
}
