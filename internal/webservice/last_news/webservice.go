package last_news

import (
	"encoding/json"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	cdb "github.com/indikator/aggregator_lets_go/internal/db/common"
	log "github.com/indikator/aggregator_lets_go/internal/log"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	wsconfig "github.com/indikator/aggregator_lets_go/internal/webservice/config"
	"net/http"
	"os"
)

type webService struct {
	config  wsconfig.Config
	log     log.Log
	db      db.Db
	handles map[string]http.Handler
}

func NewWebservice() webservice.Webservice {
	return &webService{}
}

func (ws *webService) Port() uint16 {
	return ws.config.Port
}

func (ws *webService) Db() db.Db {
	return ws.db
}

func (ws *webService) Logger() log.Log {
	return ws.log
}

func (ws *webService) Init(config *wsconfig.Config, l log.Log, db db.Db) {
	ws.config = *config
	ws.db = db
	ws.log = l
}

func (ws *webService) InitAllByConfig(config *config.Config) error {
	config.Read()
	l, err := clog.GetLog(config.WebService.Log)
	if err != nil {
		log.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("WebService.InitAllByConfig.Begin")

	db, err := cdb.GetDb(config.Database, l)
	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}
	db.InitDb()
	ws.Init(&config.WebService, l, db)

	l.WriteInfo("WebService.InitAllByConfig.End")

	return nil
}

// GetLastNews godoc
// @Summary Retrieves news from last 7 days
// @Description Get array of Articles for the last 7 days
// @Tags news
// @Produce json
// @Success 200 {object} []model.DBArticle
// @Router /last_news [get]
func (ws *webService) GetLastNews(nDays int) http.Handler {
	ws.log.WriteInfo("WebService.GetLastNews.Begin")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		news, _ := ws.db.ReadArticles(nDays)
		newsJson, _ := json.Marshal(news)
		ws.log.WriteInfo("WebService.GetLastNews.End")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		w.Write(newsJson)
	})
}

// GetBuildVersion godoc
// @Summary Retrieves build version
// @Description Get the version of current build
// @Tags Build Version
// @Produce text/plain
// @Success 200
// @Router / [get]
func (ws *webService) GetBuildVersion() http.Handler {
	ws.log.WriteInfo("WebService.GetBuildVersion.Begin")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buildVersion := os.Getenv("TAG")
		ws.log.WriteInfo("WebService.GetBuildVersion.End")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(buildVersion))
	})
}
