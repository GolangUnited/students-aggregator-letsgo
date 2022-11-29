package main

import (
	"fmt"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"strconv"

	_ "github.com/indikator/aggregator_lets_go/cmd/webservice/docs"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
	_ "github.com/indikator/aggregator_lets_go/model"
)

// @title Web-Service Swagger
// @version 1.0.1
// @description Swagger Web-Service for Let's Go Aggregator
// @termsOfService http://swagger.io/terms/

// @contact.name Web-Service Support
// @contact.email aggregator_lets_go@gmail.com

//@host https://letsgo.iost.at:8080

const (
	configFilePath = "./configs/config.yaml"
)

func main() {
	// set and read config
	cfg := config.NewConfig()
	if err := cfg.SetDataFromFile(configFilePath); err != nil {
		log.Fatal(err)
	}
	if err := cfg.Read(); err != nil {
		log.Fatal(err)
	}

	// set port from config
	port := int(cfg.WebService["last_news"].Port)

	// init db
	db := mongo.NewDb(cfg.Database)
	if err := db.DBInit(); err != nil {
		log.Fatal(err)
	}

	// init webservice
	ws := last_news.NewWebservice(cfg.WebService)

	// add swagger route to multiplexer
	r := mux.NewRouter()
	swaggerUrl := fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerUrl),
		httpSwagger.DeepLinking(true),
	)).Methods(http.MethodGet)

	// init last news handle
	lastNewsHandler := ws.GetLastNews(db, 7) // hardcode 1 week

	// create map for handles
	handlers := map[string]*http.Handler{
		"last_news": &lastNewsHandler, // add new handles if needed
	}

	// add handles to multiplexer
	for name, handle := range cfg.WebService {
		r.Handle(handle.Handle, *handlers[name])
	}

	// init and run server on given port
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
