package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/indikator/aggregator_lets_go/internal/db"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"

	_ "github.com/indikator/aggregator_lets_go/cmd/webservice/docs"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
	_ "github.com/indikator/aggregator_lets_go/model"
)

// @title Swagger Web-Service
// @version 1.0
// @description Swagger Web-Service for Let's Go Aggregator
// @termsOfService http://swagger.io/terms/

// @contact.name Web-Service Support
// @contact.email mishajudoist@gmail.com

const (
	configFilePath = "./configs/config.yaml"
)

// GetLastNews godoc
// @Summary Retrieves last news
// @Produce json
// @Success 200 {object} []model.DBArticle
// @Router /last_news [get]
func GetLastNews(db db.Db) http.Handler {

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

func main() {

	cfg := config.NewConfig()
	if err := cfg.SetDataFromFile(configFilePath); err != nil {
		log.Fatal(err)
	}

	if err := cfg.Read(); err != nil {
		log.Fatal(err)
	}

	db := mongo.NewDb(cfg.Database)
	if err := db.DBInit(); err != nil {
		log.Fatal(err)
	}

	ws := last_news.NewWebservice(cfg.WebService)

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		//httpSwagger.DocExpansion("full"),
		//httpSwagger.DomID("#swagger-ui"),
	)).Methods(http.MethodGet)

	lastNewsHandler := GetLastNews(db)
	handlers := map[string]*http.Handler{
		"last_news": &lastNewsHandler,
	}

	ws.RunServer(r, handlers)
}
