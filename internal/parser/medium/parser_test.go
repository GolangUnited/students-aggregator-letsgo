package medium_test

import (
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/parser/medium"
)

const (
	URL = "https://medium.com/_/graphql"
)

func TestParseAll(t *testing.T) {

	parser := medium.NewParser(URL)
	articles, err := parser.ParseAll()

	if err != nil {
		t.Errorf("articlesparser.ParseAll() unexpected error = %v", err)
		return
	}

	if len(articles) == 0 {
		t.Errorf("articlesparser.ParseAll() responce is empty")
		return
	}
}
