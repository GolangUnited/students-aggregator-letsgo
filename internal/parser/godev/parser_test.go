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
	news, _ := parser.ParseAll()

	if len(news) == 0 {
		t.Errorf("error: no any news %d\n", len(news))
	}
}

func TestParseBefore(t *testing.T) {

	date := time.Now()

	parser := godev.NewParser(URL)
	posts, err := parser.ParseBefore(date)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	for _, post := range posts {
		if !post.Date.Before(date) {
			t.Errorf("error: a post date %v isn't before %v\n", post.Date, date)
		}
	}
}

func TestParseBeforeN(t *testing.T) {

	const maxPostsQuantity = 3

	parser := godev.NewParser(URL)
	posts, err := parser.ParseBeforeN(time.Now(), maxPostsQuantity)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	if len(posts) != maxPostsQuantity {
		t.Errorf("error: amount of parsed posts doesn't equal to max quantity value, %d != %d \n", len(posts), maxPostsQuantity)
	}
}
