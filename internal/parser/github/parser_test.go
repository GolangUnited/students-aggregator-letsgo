package github

import (
	"testing"
	"time"
)

const (
	URL = "https://github.com/golang/go/tags"
)

func TestParseAll(t *testing.T) {

	parser := NewParser(URL)
	articles, _ := parser.ParseAll()

	if len(articles) == 0 {
		t.Errorf("error: no any news %d\n", len(articles))
	}
}

func TestParseAfter(t *testing.T) {

	date := time.Now()

	parser := NewParser(URL)
	articles, err := parser.ParseAfter(date)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	for _, article := range articles {
		if !article.Created.Before(date) {
			t.Errorf("error: a post date %v isn't before %v\n", article.Created, date)
		}
	}
}

func TestParseAfterN(t *testing.T) {

	const maxPostsQuantity = 3

	parser := NewParser(URL)
	articles, err := parser.ParseAfterN(time.Now(), maxPostsQuantity)

	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	if len(articles) != maxPostsQuantity {
		t.Errorf("error: amount of parsed posts doesn't equal to max quantity value, %d != %d \n", len(articles), maxPostsQuantity)
	}
}
