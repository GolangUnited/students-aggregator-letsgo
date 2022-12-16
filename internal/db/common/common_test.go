package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/db"
	sdb "github.com/indikator/aggregator_lets_go/internal/db/stub"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	slog "github.com/indikator/aggregator_lets_go/internal/log/stub"
)

func TestGetDbCorrectDb(t *testing.T) {
	c := db.Config{
		Name: "stub",
		Url:  "stub://localhost:22222/",
	}
	l := slog.NewLog(logLevel.Errors)
	d, err := GetDb(c, l)

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
	_, err := GetDb(c, l)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unknownDbmsError *UnknownDbmsError

	switch {
	case errors.As(err, &unknownDbmsError):
	default:
		t.Errorf("unexpected error %v", err)
	}

	got := err.Error()
	expected := fmt.Sprintf(unknownDbmsErrorTemplate, c.Name)

	if got != expected {
		t.Errorf("error message \"%s\" incorrect expected \"%s\"", got, expected)
	}
}

func TestGetDbIncorrectUrl(t *testing.T) {
	c := db.Config{
		Name: "stub",
		Url:  "",
	}
	l := slog.NewLog(logLevel.Errors)
	_, err := GetDb(c, l)

	// error: incorrect db url
	if err == nil {
		t.Error("expect error is missing")
	}
}

func TestGetDbMongo(t *testing.T) {
	originNewDb := mongoNewDb

	defer func() {
		mongoNewDb = originNewDb
	}()

	mongoNewDb = sdb.NewDb

	c := db.Config{
		Name: "mongo",
		Url:  "mongo://localhost:22222/",
	}
	l := slog.NewLog(logLevel.Errors)
	_, err := GetDb(c, l)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
