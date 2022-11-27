package aggregator

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	cdb "github.com/indikator/aggregator_lets_go/internal/db/common"
	sdb "github.com/indikator/aggregator_lets_go/internal/db/stub"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	slog "github.com/indikator/aggregator_lets_go/internal/log/stub"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	configFilePath                      = "../../tests/configs/aggregator/config.yaml"
	incorrectLogConfigFilePath          = "../../tests/configs/aggregator/incorrectLogConfig.yaml"
	incorrectParsersConfigFilePath      = "../../tests/configs/aggregator/incorrectParsersConfig.yaml"
	incorrectDbConfigFilePath           = "../../tests/configs/aggregator/incorrectDbConfig.yaml"
	incorrectEmptyParsersConfigFilePath = "../../tests/configs/aggregator/incorrectEmptyParsersConfig.yaml"
	incorrectParserUrlConfigFilePath    = "../../tests/configs/aggregator/incorrectParserUrlConfig.yaml"
	incorrectDbUrlConfigFilePath        = "../../tests/configs/aggregator/incorrectDbUrlConfig.yaml"
	daysAgoFromDb                       = 3
	incorrectParserName                 = "mock"
)

func TestIncorrectConfig(t *testing.T) {
	c := config.NewConfig()

	a := NewAggregator()

	err := a.InitAllByConfig(c)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}
func TestIncorrectLogConfig(t *testing.T) {
	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectLogConfigFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}
func TestIncorrectDbConfig(t *testing.T) {
	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectDbConfigFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}
func TestIncorrectParsersConfig(t *testing.T) {
	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectParsersConfigFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}

func TestIncorrectEmptyParsersConfig(t *testing.T) {
	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectEmptyParsersConfigFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}

func TestIncorrectInit(t *testing.T) {

	c := config.NewConfig()

	err := c.SetDataFromFile(configFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	l, err := clog.GetLog(c.Aggregator.Log)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	parsers, err := GetParsers(c.Parsers, l)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	db := sdb.NewDb(c.Database, l)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.Init(nil, l, parsers, db)

	if err == nil {
		t.Error("no error, but it was expected")
	}

	err = a.Init(&c.Aggregator, nil, parsers, db)

	if err == nil {
		t.Error("no error, but it was expected")
	}

	err = a.Init(&c.Aggregator, l, nil, db)

	if err == nil {
		t.Error("no error, but it was expected")
	}

	err = a.Init(&c.Aggregator, l, parsers, nil)

	if err == nil {
		t.Error("no error, but it was expected")
	}
}

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

func TestErrorInParseAfter(t *testing.T) {

	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectParserUrlConfigFilePath)

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

	if err == nil {
		t.Error("no error, but it was expected")
	}
}

func TestErrorInDbAfter(t *testing.T) {

	c := config.NewConfig()

	err := c.SetDataFromFile(incorrectDbUrlConfigFilePath)

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

	if err == nil {
		t.Error("no error, but it was expected")
	}
}

func TestGetParsersCorrectParser(t *testing.T) {
	pc := []parser.Config{{
		Name:    "stub",
		URL:     "https://stub.com",
		IsLocal: false,
	}}
	l := slog.NewLog(logLevel.Errors)
	p, err := GetParsers(pc, l)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(p) != len(pc) {
		t.Errorf("parsers count %d not equals parser configs count %d", len(p), len(pc))
	}
}

func TestGetParsersIncorrectParser(t *testing.T) {
	pc := []parser.Config{{
		Name:    incorrectParserName,
		URL:     "https://mock.com",
		IsLocal: false,
	}}
	l := slog.NewLog(logLevel.Errors)
	_, err := GetParsers(pc, l)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unsupportedParserNameError *UnsupportedParserNameError

	switch {
	case errors.As(err, &unsupportedParserNameError):
		em := fmt.Sprintf(unsupportedParserNameErrorTemplate, incorrectParserName)
		if err.Error() != em {
			t.Errorf("error message \"%s\" incorrect expected \"%s\"", err.Error(), em)
		}
	default:
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetDbCorrectDb(t *testing.T) {
	c := db.Config{
		Name: "stub",
		Url:  "stub://localhost:22222/",
	}
	l := slog.NewLog(logLevel.Errors)
	d, err := cdb.GetDb(c, l)

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
	l := slog.NewLog(logLevel.Errors)
	_, err := cdb.GetDb(c, l)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unknownDbmsError *cdb.UnknownDbmsError

	switch {
	case errors.As(err, &unknownDbmsError):
	default:
		t.Errorf("unexpected error %v", err)
	}
}
