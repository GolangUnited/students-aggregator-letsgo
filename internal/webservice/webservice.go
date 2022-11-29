package webservice

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"net/http"
)

type Webservice interface {
	GetLastNews(db db.Db, nDays int) http.Handler
}
