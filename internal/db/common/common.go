package common

import (
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/db/stub"
	"github.com/indikator/aggregator_lets_go/internal/log"
)

const (
	unknownDbmsErrorTemplate = "unknown dbms %s"
)

type UnknownDbmsError struct {
	Text string
}

func (e *UnknownDbmsError) Error() string {
	return e.Text
}

var (
	mongoNewDb = mongo.NewDb
)

func GetDb(config db.Config, l log.Log) (db.Db, error) {
	var d db.Db

	switch config.Name {
	case "stub":
		d = stub.NewDb(config, l)
	case "mongo":
		mongo.DecodeResults_ = mongo.DecodeResults
		mongo.CloseCursor_ = mongo.CloseCursor
		mongo.CreateClient_ = mongo.CreateClient
		mongo.PingClient_ = mongo.PingClient
		mongo.CreateDatabase_ = mongo.CreateDatabase
		mongo.CreateIndex_ = mongo.CreateIndex
		d = mongoNewDb(config, l)
	default:
		return nil, &UnknownDbmsError{
			Text: fmt.Sprintf(unknownDbmsErrorTemplate, config.Name),
		}
	}

	err := d.InitDb()

	if err != nil {
		return nil, err
	}

	return d, nil
}
