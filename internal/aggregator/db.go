package aggregator

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/db/stub"
)

func GetDb(config db.Config) (db.Db, error) {
	var d db.Db

	switch config.Name {
	case "stub":
		d = stub.NewDb(config.Url)
	case "mongo":
		d = mongo.NewDb(config.Url)
	default:
		return nil, fmt.Errorf("unknown dbms %s", config.Name)
	}

	d.DBInit(config.Url)

	return d, nil
}
