package godev_test

import (
	"errors"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	log "github.com/indikator/aggregator_lets_go/internal/log/stub"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/parser/godev"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	URL                  = "file://./page.html"
	noHref               = "file://./page-article-without-url.html"
	noTitle              = "file://./page-article-without-title.html"
	noDescription        = "file://./page-article-without-description.html"
	noDatetime           = "file://./page-article-without-datetime.html"
	noAuthor             = "file://./page-article-without-author.html"
	webPageNotFound      = "file://./page-not-found.html"
	dateFormat           = "2 January 2006"
	stringDate           = "2 August 2022"
	targetArticlesAmount = 3
)

func TestParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: URL, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	articles, err := articleParser.ParseAfter(date)
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

func TestWithoutTitleParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: noTitle, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorArticleTitleNotFound) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestWithoutAuthorParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: noAuthor, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorArticleAuthorNotFound) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestWithoutDescriptionParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: noDescription, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorArticleDescriptionNotFound) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestWithoutHrefParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: noHref, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorArticleURLNotFound) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestWithoutDatetimeParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: noDatetime, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorArticleDatetimeNotFound) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func TestWebPageNotFoundParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: webPageNotFound, IsLocal: true}, log.NewLog(logLevel.Errors)

	articleParser := godev.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}

	_, err = articleParser.ParseAfter(date)
	if err != nil {
		if !errors.Is(err, parser.ErrorWebPageCannotBeDelivered{URL: webPageNotFound, StatusCode: 404}) {
			t.Errorf("error: %s\n", err)
		}
	}
}

func getArticles() (articles []model.Article, err error) {

	datestrings := []string{"26 September 2022", "8 September 2022", "6 September 2022"}
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
			Title:       "Go runtime: 4 years later",
			Created:     dates[0],
			Description: "A check-in on the status of Go runtime development",
			URL:         "file://./blog/go119runtime",
			Author:      "Michael Knyszek",
		}, {
			Title:       "Go Developer Survey 2022 Q2 Results",
			Created:     dates[1],
			Description: "An analysis of the results from the 2022 Q2 Go Developer Survey.",
			URL:         "file://./blog/survey2022-q2-results",
			Author:      "Todd Kulesza",
		}, {
			Title:       "Vulnerability Management for Go",
			Created:     dates[2],
			Description: "Announcing vulnerability management for Go, to help developers learn about known vulnerabilities in their dependencies.",
			URL:         "file://./blog/vuln",
			Author:      "Julie Qiu, for the Go security team",
		},
	}, nil
}
