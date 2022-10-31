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
		news, err := json.Marshal(news)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(news)
	})
}

func RunServer(ws webservice.Webservice, cfg config.Config, handle string) error {
	if err := cfg.Read(); err != nil {
		return err
	}

	db := mongo.NewDb(cfg.Database.Url)
	if err := db.DBInit(); err != nil {
		fmt.Println("BBB")
		return fmt.Errorf("can't start the server: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle(handle, ws.MessageHandler(db))

	log.Println("Listening... on port: ", cfg.WebService.Port)
	http.ListenAndServe(":"+strconv.Itoa(int(cfg.WebService.Port)), mux)

	return nil
}
