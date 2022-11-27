package aggregator

import (
	"fmt"
	"time"

	aconfig "github.com/indikator/aggregator_lets_go/internal/aggregator/config"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	cdb "github.com/indikator/aggregator_lets_go/internal/db/common"
	"github.com/indikator/aggregator_lets_go/internal/log"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Aggregator struct {
	config            aconfig.Config
	log               log.Log
	parsers           []parser.ArticlesParser
	db                db.Db
	lastCheckDatetime time.Time
}

const (
	daysAgo = 1
)

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) InitAllByConfig(config *config.Config) error {
	err := config.Read()

	if err != nil {
		log.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	l, err := clog.GetLog(config.Aggregator.Log)

	if err != nil {
		log.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("Aggregator.InitAllByConfig.Begin")

	parsers, err := GetParsers(config.Parsers, l)

	if err != nil {
		l.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	db, err := cdb.GetDb(config.Database, l)

	if err != nil {
		l.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	err = a.Init(&config.Aggregator, l, parsers, db)

	if err != nil {
		l.WriteError("Aggregator.InitAllByConfig.Error", err)
		return err
	}

	l.WriteInfo("Aggregator.InitAllByConfig.End")

	return nil
}

func (a *Aggregator) Init(config *aconfig.Config, l log.Log, parsers []parser.ArticlesParser, db db.Db) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if l == nil {
		return fmt.Errorf("log is nil")
	}
	if parsers == nil {
		return fmt.Errorf("parsers is nil")
	}
	if db == nil {
		return fmt.Errorf("db is nil")
	}
	a.config = *config
	a.log = l
	a.parsers = parsers
	a.db = db

	return nil
}

func composeError(errAll, err error) error {
	if errAll != nil {
		errAll = fmt.Errorf("%v; %v", errAll, err)
	} else {
		errAll = err
	}
	return errAll
}

func (a *Aggregator) Execute() error {
	a.log.WriteInfo("Aggregator.Execute.Begin")

	a.lastCheckDatetime = time.Now().AddDate(0, -daysAgo, 0)

	var errAll error = nil

	for _, v := range a.parsers {
		articles, err := v.ParseAfter(a.lastCheckDatetime)

		if err != nil {
			a.log.WriteError("Aggregator.Execute.Error", err)
			errAll = composeError(errAll, err)
			continue
		}

		for _, article := range articles {
			id := primitive.NewObjectID()

			_, err = a.db.WriteArticle(&model.DBArticle{
				ID:          id,
				Title:       article.Title,
				Created:     article.Created,
				Author:      article.Author,
				Description: article.Description,
				URL:         article.URL,
			})

			if err != nil {
				a.log.WriteError("Aggregator.Execute.Error", err)
				errAll = composeError(errAll, err)
			}
		}
	}

	a.log.WriteInfo("Aggregator.Execute.End")

	return errAll
}
