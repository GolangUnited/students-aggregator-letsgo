package webservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func MessageHandler(news []DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		news, _ := json.Marshal(news)
		w.Write(news)
	})
}

func RunServer(news []DB) error {
	mux := http.NewServeMux()
	mux.Handle("/last_news", MessageHandler(news))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		return fmt.Errorf("Can't start the server")
	}
	log.Println("Listening...")
	return nil
}
