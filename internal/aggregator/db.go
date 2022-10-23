package aggregator

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/db/stub"
)

func GetDb(config db.Config) (db.Db, error) {
	if config.Name == "Stub" {
		d := stub.Db{}

		d.DBinit(config.Url)

		return &d, nil
	}

	return nil, fmt.Errorf("unknown dbms %s", config.Name)
}
