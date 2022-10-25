package webservice

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"net/http"
)

type Webservice interface {
	MessageHandler(db db.Db) http.Handler
}
