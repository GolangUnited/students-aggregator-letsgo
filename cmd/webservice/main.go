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

	// set port from config
	port := int(ws.Port())

	// add swagger route to multiplexer
	r := mux.NewRouter()
	swaggerUrl := fmt.Sprintf("http://localhost:%d/swagger/doc.json", port)
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL(swaggerUrl),
		httpSwagger.DeepLinking(true),
	)).Methods(http.MethodGet)

	// init last news handle
	lastNewsHandler := ws.GetLastNews(7) // hardcode 1 week
	r.Handle(cfg.WebService.Handle, lastNewsHandler)

	// init and run server on given port
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
