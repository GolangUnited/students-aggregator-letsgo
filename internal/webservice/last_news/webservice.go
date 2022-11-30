package last_news

import (
	"encoding/json"
	"fmt"
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

func (ws *webService) Db() *db.Db {
	return &ws.db
}

func (ws *webService) Logger() *log.Log {
	return &ws.log
}

func (ws *webService) Init(config *wsconfig.Config, l log.Log, db db.Db) error {
	ws.config = *config
	ws.db = db
	ws.log = l
	return nil
}

func (ws *webService) InitAllByConfig(config *config.Config) error {
	err := config.Read()

	l, err := clog.GetLog(config.WebService.Log)

	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("WebService.InitAllByConfig.Begin")

	db, err := cdb.GetDb(config.Database, l)

	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	err = ws.Init(&config.WebService, l, db)

	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("WebService.InitAllByConfig.End")

	return nil
}

func LoggingHandler(next http.Handler, l *log.Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(fmt.Sprintf("User %v using app", r.RemoteAddr))
		(*l).WriteInfo(fmt.Sprintf("User %v using app", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

// GetLastNews godoc
// @Summary Retrieves news from last 7 days
// @Description Get array of Articles for the last 7 days
// @Tags news
// @Produce json
// @Success 200 {object} []model.DBArticle
// @Router /last_news [get]
func (ws *webService) GetLastNews(nDays int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.log.WriteInfo("WebService.GetLastNews.Begin")
		news, err := ws.db.ReadArticles(nDays)
		if err != nil {
			ws.log.WriteError("WebService.GetLastNews.Error", err)
			os.Exit(1)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		newsJson, err := json.Marshal(news)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(newsJson)
		ws.log.WriteInfo("WebService.GetLastNews.End")
	})
}
