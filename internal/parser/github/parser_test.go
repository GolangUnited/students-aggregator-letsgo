package github_test

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/parser/github"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	URL                  = "file://./page_github_tags.html"
	dateFormat           = "2006-01-02T15:04:05Z"
	stringDate           = "2022-10-04T00:00:00Z"
	targetArticlesAmount = 4
)

func TestParseAfter(t *testing.T) {
	cfg := parser.Config{URL: URL, IsLocal: true}

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	parser := github.NewParser(cfg)
	articles, err := parser.ParseAfter(date)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	if len(articles) != targetArticlesAmount {
		t.Errorf("error: wromg amount of articles (%d != %d) after %v", targetArticlesAmount, len(articles), date)
	}

	testesArticles, err := getArticles()
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	for index, article := range articles {
		if !article.Created.After(date) {
			t.Errorf("error: an article date %v isn't after %v\n", article.Created, date)
		}
		if !article.Created.Equal(testesArticles[index].Created) {
			t.Errorf("error: an parsed article date: %v is not equal to the expected article date: %v\n", article.Created, testesArticles[index].Created)
		}
		if article.Title != testesArticles[index].Title {
			t.Errorf("error: an parsed article title: %s is not equal to the expected article title: %s\n", article.Title, testesArticles[index].Title)
		}
		if article.Description != testesArticles[index].Description {
			t.Errorf("error: an parsed article description: %s is not equal to the expected article description: %s\n", article.Description, testesArticles[index].Description)
		}
		if article.URL != testesArticles[index].URL {
			t.Errorf("error: an parsed article URL: %s is not equal to the expected article URL: %s\n", article.URL, testesArticles[index].URL)
		}
		if article.Author != testesArticles[index].Author {
			t.Errorf("error: an parsed article Author: %s is not equal to the expected article Author: %s\n", article.Author, testesArticles[index].Author)
		}
	}
}

func getArticles() (articles []model.Article, err error) {

	datestrings := []string{"2022-11-01T16:45:23Z", "2022-11-01T16:45:18Z", "2022-10-04T17:43:19Z", "2022-10-04T17:43:09Z"}
	dates := make([]time.Time, 0)
	for _, datestring := range datestrings {
		date, err := time.Parse(dateFormat, datestring)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}

	return []model.Article{
		{
			Title:       "go1.19.3",
			URL:         "https://github.com/golang/go/releases/tag/go1.19.3",
			Created:     dates[0],
			Author:      "Gopher Robot <gobot@golang.org>",
			Description: "[release-branch.go1.19] go1.19.3\n      Change-Id: I167308920eeb7480efb626ce75f777a335e870b0\nReviewed-on: https://go-review.googlesource.com/c/go/+/446958\nRun-TryBot: Gopher Robot <gobot@golang.org>\nReviewed-by: Matthew Dempsky <mdempsky@google.com>\nAuto-Submit: Gopher Robot <gobot@golang.org>\nReviewed-by: Heschi Kreinick <heschi@google.com>\nTryBot-Result: Gopher Robot <gobot@golang.org>",
		}, {
			Title:       "go1.18.8",
			URL:         "https://github.com/golang/go/releases/tag/go1.18.8",
			Created:     dates[1],
			Author:      "Gopher Robot <gobot@golang.org>",
			Description: "[release-branch.go1.18] go1.18.8\n      Change-Id: I89e791f1d6ae0984ba62bccef05886acbb10b2dd\nReviewed-on: https://go-review.googlesource.com/c/go/+/446957\nRun-TryBot: Gopher Robot <gobot@golang.org>\nReviewed-by: Matthew Dempsky <mdempsky@google.com>\nTryBot-Result: Gopher Robot <gobot@golang.org>\nAuto-Submit: Gopher Robot <gobot@golang.org>\nReviewed-by: Heschi Kreinick <heschi@google.com>",
		}, {
			Title:       "go1.19.2",
			URL:         "https://github.com/golang/go/releases/tag/go1.19.2",
			Created:     dates[2],
			Author:      "Gopher Robot <gobot@golang.org>",
			Description: "[release-branch.go1.19] go1.19.2\n      Change-Id: Ia5de3a0fa07f212c5c19f9e01b0ed2cfab739e95\nReviewed-on: https://go-review.googlesource.com/c/go/+/438598\nReviewed-by: Dmitri Shuralyov <dmitshur@google.com>\nReviewed-by: Carlos Amedee <carlos@golang.org>\nAuto-Submit: Gopher Robot <gobot@golang.org>\nRun-TryBot: Gopher Robot <gobot@golang.org>\nTryBot-Result: Gopher Robot <gobot@golang.org>",
		}, {
			Title:       "go1.18.7",
			URL:         "https://github.com/golang/go/releases/tag/go1.18.7",
			Created:     dates[3],
			Author:      "Gopher Robot <gobot@golang.org>",
			Description: "[release-branch.go1.18] go1.18.7\n      Change-Id: I0636d7335381c25ce39fd44c8cf758fb84737551\nReviewed-on: https://go-review.googlesource.com/c/go/+/438597\nReviewed-by: Carlos Amedee <carlos@golang.org>\nRun-TryBot: Gopher Robot <gobot@golang.org>\nReviewed-by: Dmitri Shuralyov <dmitshur@google.com>\nAuto-Submit: Gopher Robot <gobot@golang.org>\nTryBot-Result: Gopher Robot <gobot@golang.org>",
		},
	}, nil
}
