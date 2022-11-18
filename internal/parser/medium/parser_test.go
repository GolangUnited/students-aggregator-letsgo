package medium_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/parser/medium"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	// The file contains a request to the service database
	url                  = "http://localhost/testMediumservice"
	responceFile         = "../../../tests/data/parser/medium/response.json"
	testsArticlesFile    = "../../../tests/data/parser/medium/testArticles.json"
	queryFileName        = "../../../etc/queryInMedium"
	dateFormat           = "2006-01-02T15:04:05Z"
	stringDate           = "2022-11-10T00:00:00Z"
	hostService          = "localhost"
	targetArticlesAmount = 5
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
	}
}

func TestParseAfter(t *testing.T) {

	tests := []struct {
		name              string
		url               string
		responceFile      string
		testsArticlesFile string
		queryFileName     string
		hostService       string
		wantErr           bool
	}{
		{name: "Default case", url: url, responceFile: responceFile, testsArticlesFile: testsArticlesFile,
			queryFileName: queryFileName, hostService: hostService, wantErr: false},
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

			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				equals(t, req.URL.String(), tt.url)
				responceBody, _ := os.ReadFile(tt.responceFile)
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBuffer(responceBody)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})
			cfg := parser.Config{URL: tt.url, IsLocal: true}

			parser := func(cfg parser.Config) parser.ArticlesParser {

				return &medium.ArticlesParser{
					Client:        client,
					Url:           cfg.URL,
					Host:          tt.hostService,
					QueryFileName: tt.queryFileName,
				}
			}(cfg)

			gotArticles, err := parser.ParseAfter(date)

			if (err != nil) != tt.wantErr {
				t.Errorf("articlesparser.ParseAfter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotArticles) != targetArticlesAmount {
				t.Errorf("error: wromg amount of articles (%d != %d) after %v", targetArticlesAmount, len(gotArticles), date)
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

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func getArticles(filename string) (articles []model.Article, err error) {

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &articles)
	if err != nil {
		return nil, err
	}

	return articles, nil

}
