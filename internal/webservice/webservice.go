package webservice

import (
	"net/http"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/log"
	wsconfig "github.com/indikator/aggregator_lets_go/internal/webservice/config"
)

type Webservice interface {
	InitAllByConfig(config *config.Config) error
	Init(config *wsconfig.Config, l log.Log, db db.Db) error
	MessageHandler() http.Handler
	RunServer()
}
