package aggregator

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/config"
)

const (
	configData = `# Project Aggregator YAML
aggregator:
  nothing:

database:
  name: stub
  url: stub://localhost:22222/

webservice:
  port: 8080

parsers:
- stub:
    url: https://stub.com`
)

func TestWorkWithStubParser(t *testing.T) {

	c := config.NewConfig()

	err := c.SetData([]byte(configData))

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(c.Parsers) != 1 {
		t.Errorf("incorrect config parsers count %d, expected %d", len(c.Parsers), 1)
	}

	parsers, err := GetParsers(c.Parsers)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(parsers) != 1 {
		t.Errorf("incorrect parsers count %d, expected %d", len(parsers), 1)
	}

	parser := parsers[0]

	date := time.Now().AddDate(0, -2, 0)
	articles, err := parser.ParseAfter(date)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(articles) != 3 {
		t.Errorf("incorrect articles count %d, expected %d", len(articles), 3)
	}

	db, err := GetDb(c.Database)

	// if db.As(db_stub.Db) {
	// 	t.Errorf("incorrect dbms, expected stub")
	// }

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.Init(&c.Aggregator, parsers, db)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = a.Execute()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	dbArticles, err := db.ReadAllArticles()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(dbArticles) != 3 {
		t.Errorf("incorrect db articles count %d, expected %d", len(dbArticles), 3)
	}
}
