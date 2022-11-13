package aggregator

import (
	"errors"
	"reflect"
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	configFilePath = "../../tests/configs/aggregator/config.yaml"
	daysAgoFromDb  = 3
)

func TestWorkWithStubParserAndDb(t *testing.T) {

	c := config.NewConfig()

	err := c.SetDataFromFile(configFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = a.Execute()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	articlesFromParser, err := a.parsers[0].ParseAfter(a.lastCheckDatetime)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	articles, err := a.db.ReadArticles(daysAgoFromDb)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	var articlesFromDb []model.Article

	for _, article := range articles {
		dbArticle := db.ConvertFromDbArticle(article)

		articlesFromDb = append(articlesFromDb, *dbArticle)
	}

	if !reflect.DeepEqual(articlesFromDb, articlesFromParser) {
		t.Log(articlesFromDb, articlesFromParser)

		t.Errorf("articles received from a database are different from articles received from a parser")
	}
}

func TestGetParsersCorrectParser(t *testing.T) {
	pc := []parser.Config{{
		Name:    "stub",
		URL:     "https://stub.com",
		IsLocal: false,
	}}

	p, err := GetParsers(pc)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(p) != len(pc) {
		t.Errorf("parsers count %d not equals parser configs count %d", len(p), len(pc))
	}
}

func TestGetParsersIncorrectParser(t *testing.T) {
	pc := []parser.Config{{
		Name:    "mock",
		URL:     "https://mock.com",
		IsLocal: false,
	}}

	_, err := GetParsers(pc)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unsupportedParserNameError *UnsupportedParserNameError

	switch {
	case errors.As(err, &unsupportedParserNameError):
	default:
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetDbCorrectDb(t *testing.T) {
	c := db.Config{
		Name: "stub",
		Url:  "stub://localhost:22222/",
	}

	d, err := GetDb(c)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if c.Name != d.Name() {
		t.Errorf("dbms name %s not equals dbms config name %s", d.Name(), c.Name)
	}

	if c.Url != d.Url() {
		t.Errorf("dbms url %s not equals dbms config url %s", d.Url(), c.Url)
	}
}

func TestGetDbIncorrectDb(t *testing.T) {
	c := db.Config{
		Name: "mock",
		Url:  "mock://localhost:22222/",
	}

	_, err := GetDb(c)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unknownDbmsError *UnknownDbmsError

	switch {
	case errors.As(err, &unknownDbmsError):
	default:
		t.Errorf("unexpected error %v", err)
	}
}
