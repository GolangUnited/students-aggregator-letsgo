package last_news

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	cdb "github.com/indikator/aggregator_lets_go/internal/db/common"
	log2 "github.com/indikator/aggregator_lets_go/internal/log"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	wsconfig "github.com/indikator/aggregator_lets_go/internal/webservice/config"
)

type webService struct {
	config wsconfig.Config
	log    log2.Log
	db     db.Db
}

func NewWebservice() webservice.Webservice {
	return &webService{}
}

func (ws *webService) InitAllByConfig(config *config.Config) error {
	err := config.Read()

	if err != nil {
		log2.WriteError("WebService.InitAllByConfig.Error", err)
		return err
	}

	l, err := clog.GetLog(config.WebService.Log)

	if err != nil {
		log2.WriteError("WebService.InitAllByConfig.Error", err)
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
		l.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("Aggregator.InitAllByConfig.End")

	return nil
}

func (ws *webService) Init(config *wsconfig.Config, l log2.Log, db db.Db) error {
	ws.config = *config
	ws.log = l
	ws.db = db

	return nil
}

func (ws *webService) MessageHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		news, err := ws.db.ReadArticles(7) // hardcode 1 week
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
	})
}

func (ws *webService) RunServer() {
	mux := http.NewServeMux()
	mux.Handle(ws.config.Handle, ws.MessageHandler())
	ws.log.WriteInfo("Listening...")
	http.ListenAndServe(":"+strconv.Itoa(int(ws.config.Port)), mux)
}
