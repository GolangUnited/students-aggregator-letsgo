package godev_test

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser/godev"
)

const (
	URL = "https://go.dev/blog/all"
)

func TestParseAfter(t *testing.T) {

	date := time.Now().AddDate(-1, 0, 0)

	parser := godev.NewParser(URL)
	articles, err := parser.ParseAfter(date)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	for _, article := range articles {
		if !article.Created.After(date) {
			t.Errorf("error: an article date %v isn't after %v\n", article.Created, date)
		}
	}
}
