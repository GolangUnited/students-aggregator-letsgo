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
		d = stub.NewDb(config)
	case "mongo":
		d = mongo.NewDb(config)
	default:
		return nil, fmt.Errorf("unknown dbms %s", config.Name)
	}

	err := d.DBInit()

	if err != nil {
		return nil, err
	}

	return d, nil
}
