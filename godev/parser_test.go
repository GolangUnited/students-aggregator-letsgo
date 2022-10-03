package godev_test

import (
	"testing"

	"github.com/indikator/aggregator_lets_go/godev"
)

const (
	URL = "https://go.dev/blog/"
)

func TestLoadAllNews(t *testing.T) {

	parser := godev.NewParser(URL)
	news, _ := parser.ParseAll()

	if len(news) == 0 {
		t.Errorf("no any news %d\n", len(news))
	}
}
