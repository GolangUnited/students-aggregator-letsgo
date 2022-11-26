package last_news

import (
	"github.com/gorilla/mux"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"log"
	"net/http"
	"strconv"
)

type webService struct {
	handles map[string]webservice.Config
}

func NewWebservice(config map[string]webservice.Config) webservice.Webservice {
	return &webService{
		handles: config,
	}
}

func (ws *webService) RunServer(r *mux.Router, handlers map[string]*http.Handler) {
	port := int(ws.handles["last_news"].Port)
	for name, handle := range ws.handles {
		r.Handle(handle.Handle, *handlers[name])
	}
	log.Println("Listening...")
	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
