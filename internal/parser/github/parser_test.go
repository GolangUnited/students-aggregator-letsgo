package github

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
)

const (
	URL = "https://github.com/golang/go/tags"
)

func TestParseAfter(t *testing.T) {

	cfg := parser.Config{URL: URL}
	parser, date := NewParser(cfg), time.Now()
	articles, err := parser.ParseAfter(date)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	for _, article := range articles {
		if !article.Created.After(date) {
			t.Errorf("error: a post date %v isn't before %v\n", article.Created, date)
		}
	}
}
