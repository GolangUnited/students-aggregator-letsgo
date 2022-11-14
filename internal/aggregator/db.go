package aggregator

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"github.com/indikator/aggregator_lets_go/internal/db/stub"
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

func GetDb(config db.Config) (db.Db, error) {
	var d db.Db

	switch config.Name {
	case "stub":
		d = stub.NewDb(config)
	case "mongo":
		d = mongo.NewDb(config)
	default:
		return nil, &UnknownDbmsError{
			Text: fmt.Sprintf(unknownDbmsErrorTemplate, config.Name),
		}
	}

	err := d.DBInit()

	if err != nil {
		return nil, err
	}

	return d, nil
}
