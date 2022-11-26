package webservice

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Webservice interface {
	RunServer(mux *mux.Router, handlers map[string]*http.Handler)
}
