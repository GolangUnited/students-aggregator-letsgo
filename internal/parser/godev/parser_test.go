package godev_test

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser/godev"
)

const (
	URL = "https://go.dev/blog/all"
)

func TestParseAll(t *testing.T) {

	parser := godev.NewParser(URL)
	articles, _ := parser.ParseAll()

	if len(articles) == 0 {
		t.Errorf("error: no any articles  %d\n", len(articles))
	}
}

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

func TestParseAfterN(t *testing.T) {

	const maxPostsQuantity = 3

	date := time.Now().AddDate(-1, 0, 0)

	parser := godev.NewParser(URL)
	articles, err := parser.ParseAfterN(date, maxPostsQuantity)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	if len(articles) != maxPostsQuantity {
		t.Errorf("error: amount of parsed articles doesn't equal to max quantity value, %d != %d \n", len(articles), maxPostsQuantity)
	}

	for _, article := range articles {
		if !article.Created.After(date) {
			t.Errorf("error: an article date %v isn't after %v\n", article.Created, date)
		}
	}
}
