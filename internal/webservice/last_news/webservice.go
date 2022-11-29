package last_news

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	cdb "github.com/indikator/aggregator_lets_go/internal/db/common"
	log "github.com/indikator/aggregator_lets_go/internal/log"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	wsconfig "github.com/indikator/aggregator_lets_go/internal/webservice/config"
)

type webService struct {
	handles map[string]webservice.Config
}

func (ws *webService) InitAllByConfig(config *config.Config) error {
	err := config.Read()

	if err != nil {
		log.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

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

func NewWebservice(config map[string]webservice.Config) webservice.Webservice {
	return &webService{
		handles: config,
	err = ws.Init(&config.WebService, l, db)

	if err != nil {
		l.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("WebService.InitAllByConfig.End")

	return nil
}

func (ws *webService) Init(config *wsconfig.Config, l log.Log, db db.Db) error {
	ws.config = *config
	ws.log = l
	ws.db = db

	return nil
}

// GetLastNews godoc
// @Summary Retrieves last news from last 7 days
// @Description Get array of Articles for the last 7 days
// @Tags news
// @Produce json
// @Success 200 {object} []model.DBArticle
// @Router /last_news [get]
func (ws *webService) GetLastNews(db db.Db, nDays int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.log.WriteInfo("WebService.MessageHandler.Begin")
		news, err := db.ReadArticles(nDays)
		if err != nil {
			ws.log.WriteError("WebService.MessageHandler.Error", err)
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
		ws.log.WriteInfo("WebService.MessageHandler.End")
	})
}
