package webservice

import (
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/log"
	wsconfig "github.com/indikator/aggregator_lets_go/internal/webservice/config"
	"net/http"
)

type Webservice interface {
	GetLastNews(nDays int) http.Handler
	InitAllByConfig(config *config.Config) error
	Init(config *wsconfig.Config, l log.Log, db db.Db)
	Port() uint16
	Db() db.Db
	Logger() log.Log
}
