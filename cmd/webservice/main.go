package main

import (
	"github.com/gorilla/mux"
	"github.com/indikator/aggregator_lets_go/internal/webservice/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"strconv"

	_ "github.com/indikator/aggregator_lets_go/cmd/webservice/docs"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
)

// @title Web-Service Swagger
// @version 1.0.1
// @description Swagger Web-Service for Let's Go Aggregator
// @termsOfService http://swagger.io/terms/

// @contact.name Web-Service Support
// @contact.email aggregator_lets_go@gmail.com

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

	// init webservice
	ws := last_news.NewWebservice()
	if err := ws.InitAllByConfig(cfg); err != nil {
		log.Fatal(err)
	}

	// set port and logger
	port := int(ws.Port())
	logger := ws.Logger()

	r := mux.NewRouter()

	// add swagger route to multiplexer
	swaggerUrl := "./swagger/doc.json"
	r.PathPrefix("/swagger/").Handler(
		httpSwagger.Handler(
			httpSwagger.URL(swaggerUrl),
			httpSwagger.DeepLinking(true),
		)).Methods(http.MethodGet)

	// init last news handle
	lastNewsHandler := ws.GetLastNews(7) // hardcode 1 week

	r.Handle(cfg.WebService.Handle, middleware.LoggingHandler(lastNewsHandler, logger))

	// init and run server on given port
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	logger.WriteInfo("Listening...")
	server.ListenAndServe()
}
