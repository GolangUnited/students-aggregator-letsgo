package medium_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	"github.com/indikator/aggregator_lets_go/internal/log/stub"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/parser/medium"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	// The file contains a request to the service database
	url           = "http://localhost/testMediumservice"
	queryFileName = "../../../etc/queryInMedium"
	dateFormat    = "2006-01-02T15:04:05Z"
	stringDate    = "2022-11-10T00:00:00Z"
	hostService   = "localhost"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
		Timeout:   time.Nanosecond * 1,
	}
}

func TestParseAfter(t *testing.T) {

	tests := []struct {
		name                 string
		url                  string
		responceFile         string
		testsArticlesFile    string
		hostService          string
		targetArticlesAmount int
		iterationCount       int
		statusCode           int
		checkRedirectWant    bool
		wantErr              error
	}{
		// {name: "Default case", url: "https://medium.com/_/graphql", responceFile: "../../../tests/data/parser/medium/response.json",
		// 	testsArticlesFile: "../../../tests/data/parser/medium/testArticles.json", iterationCount: 4,
		// 	hostService: "medium.com", targetArticlesAmount: 5, statusCode: 200, wantErr: nil},
		{name: "Default case", url: url, responceFile: "../../../tests/data/parser/medium/response.json",
			testsArticlesFile: "../../../tests/data/parser/medium/testArticles.json", iterationCount: 1,
			hostService: hostService, targetArticlesAmount: 5, statusCode: 200, wantErr: nil},
		{name: "Empty title", url: url, responceFile: "../../../tests/data/parser/medium/response-without-title.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: parser.ErrorArticleTitleNotFound},
		{name: "Empty url", url: url, responceFile: "../../../tests/data/parser/medium/response-without-url.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: parser.ErrorArticleURLNotFound},
		{name: "Empty author", url: url, responceFile: "../../../tests/data/parser/medium/response-without-author.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: parser.ErrorArticleAuthorNotFound},
		{name: "Empty datetime", url: url, responceFile: "../../../tests/data/parser/medium/response-without-datetime.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: parser.ErrorArticleDatetimeNotFound},
		{name: "Empty description", url: url, responceFile: "../../../tests/data/parser/medium/response-without-description.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: parser.ErrorArticleDescriptionNotFound},
		{name: "Web page not found", url: url, responceFile: "../../../tests/data/parser/medium/response-page-not-found.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 404, wantErr: parser.ErrorWebPageCannotBeDelivered{URL: url, StatusCode: 404}},
		{name: "Unknown error (status code 500)", url: url, responceFile: "../../../tests/data/parser/medium/response-without-description.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 500, wantErr: &parser.ErrorUnknown{}},
		{name: "Unknown error (wrong json unmarshall)", url: url, responceFile: "../../../tests/data/parser/medium/response-empty.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: &parser.ErrorUnknown{}},
		{name: "Unknown error (empty item)", url: url, responceFile: "../../../tests/data/parser/medium/response-empty-item.json",
			testsArticlesFile: "", hostService: hostService, targetArticlesAmount: 0,
			iterationCount: 1, statusCode: 200, wantErr: &parser.ErrorUnknown{}},
		{name: "Unknown error (bad url)", url: "bad url/\t", responceFile: "../../../tests/data/parser/medium/response.json",
			testsArticlesFile: "", iterationCount: 1, hostService: hostService, targetArticlesAmount: 0,
			statusCode: 200, wantErr: &parser.ErrorUnknown{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			testsArticles, err := getArticles(tt.testsArticlesFile)
			if err != nil {
				t.Errorf("error: %s\n", err.Error())
			}

			date, err := time.Parse(dateFormat, stringDate)
			if err != nil {
				t.Errorf("error: %s\n", err.Error())
			}

			responceBody, _ := os.ReadFile(tt.responceFile)

			client := NewTestClient(func(req *http.Request) *http.Response {

				return &http.Response{
					StatusCode: tt.statusCode,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBuffer(responceBody)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			cfg, lg := parser.Config{URL: tt.url, IsLocal: true}, stub.NewLog(logLevel.Errors)

			newParser := func(cfg parser.Config, lg log.Log) parser.ArticlesParser {

				return &medium.ArticlesParser{
					Client:         client,
					Url:            cfg.URL,
					Host:           tt.hostService,
					QueryFileName:  queryFileName,
					IterationCount: tt.iterationCount,
					LocalLaunch:    cfg.IsLocal,
					Log:            lg,
				}
			}(cfg, lg)

			gotArticles, err := newParser.ParseAfter(date)

			if !errors.Is(err, tt.wantErr) {
				if !errors.As(err, &parser.ErrorUnknown{}) {
					t.Errorf("articlesparser.ParseAfter() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if len(gotArticles) != tt.targetArticlesAmount {
				t.Errorf("error: wromg amount of articles (%d != %d) after %v", tt.targetArticlesAmount, len(gotArticles), date)
			}

			for index, article := range gotArticles {

				if !article.Created.After(date) {
					t.Errorf("error: an article date %v isn't after %v\n", article.Created, date)
				}
				if !article.Created.Equal(testsArticles[index].Created) {
					t.Errorf("error: an parsed article date: %v is not equal to the expected article date: %v\n", article.Created,
						testsArticles[index].Created)
				}
				if article.Title != testsArticles[index].Title {
					t.Errorf("error: an parsed article title: %s is not equal to the expected article title: %s\n", article.Title,
						testsArticles[index].Title)
				}
				if article.Description != testsArticles[index].Description {
					t.Errorf("error: an parsed article description: %s is not equal to the expected article description: %s\n",
						article.Description, testsArticles[index].Description)
				}
				if article.URL != testsArticles[index].URL {
					t.Errorf("error: an parsed article URL: %s is not equal to the expected article URL: %s\n", article.URL,
						testsArticles[index].URL)
				}
				if article.Author != testsArticles[index].Author {
					t.Errorf("error: an parsed article Author: %s is not equal to the expected article Author: %s\n", article.Author,
						testsArticles[index].Author)
				}
			}

		})
	}
}

func TestQueryFileNameNotFoundParseAfter(t *testing.T) {

	cfg, lg := parser.Config{URL: url, IsLocal: true}, stub.NewLog(logLevel.Errors)

	articlesParser := medium.NewParser(cfg, lg)

	date, err := time.Parse(dateFormat, stringDate)
	if err != nil {
		t.Errorf("error: %s\n", err.Error())
	}

	_, err = articlesParser.ParseAfter(date)
	if err != nil {
		if !errors.As(err, &parser.ErrorUnknown{}) {
			t.Errorf("error: %s\n", err)
		}
	} else {
		t.Error("error cannot equals to nil")
	}
}

func TestErrorRequestAndReadBodyParseAfter(t *testing.T) {

	serverMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1")
	}))
	defer serverMock.Close()

	urls := [2]string{
		serverMock.URL,
		url,
	}

	for _, value := range urls {

		cfg, lg := parser.Config{URL: value, IsLocal: true}, stub.NewLog(logLevel.Errors)

		newParser := func(cfg parser.Config, lg log.Log) parser.ArticlesParser {

			return &medium.ArticlesParser{
				Client:         &http.Client{},
				Url:            cfg.URL,
				Host:           hostService,
				QueryFileName:  queryFileName,
				IterationCount: 1,
				LocalLaunch:    cfg.IsLocal,
				Log:            lg,
			}
		}(cfg, lg)

		date, err := time.Parse(dateFormat, stringDate)
		if err != nil {
			t.Errorf("error: %s\n", err.Error())
		}

		_, err = newParser.ParseAfter(date)
		if err != nil {
			if !errors.As(err, &parser.ErrorUnknown{}) {
				t.Errorf("error: %s\n", err)
			}
		} else {
			t.Error("error cannot equals to nil")
		}
	}

}

func getArticles(filename string) (articles []model.Article, err error) {

	if len(filename) == 0 {
		return
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &articles)
	if err != nil {
		return
	}

	return

}
