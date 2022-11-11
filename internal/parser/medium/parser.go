package medium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	limitQuantityStates = 25
	// The file contains a request to the service database
	queryFileName = "./etc/queryInMedium"
)

type articlesParser struct {
	url string
}

type StatesMedium struct {
	Data DataStates
}

type DataStates struct {
	TagFeed TagFeed
}

type TagFeed struct {
	Items []Item
}

type Item struct {
	FeedId string
	Post   Post
}

type Post struct {
	Id                     string
	Creator                Creator
	ExtendedPreviewContent ExtendedPreviewContent
	FirstPublishedAt       int64
	Title                  string
	MediumUrl              string
}

type Creator struct {
	Name string
}

type ExtendedPreviewContent struct {
	Subtitle string
}

// create an instance of articles parser
func NewParser(cfg parser.Config) parser.ArticlesParser {
	return &articlesParser{
		url: cfg.URL,
	}
}

func init() {
	parser.RegisterParser("medium.com", NewParser)
}

// / parse all articles that were created earler than the target date
func (p *articlesParser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	var (
		dateLastState   int64
		initNumberState int
	)

	maxDateUnix := maxDate.UnixMilli()
	initialRequest := true

	for true {

		states, err := getStates(p.url, initialRequest, initNumberState)
		if err != nil {
			return nil, err
		}

		if len(states) == 0 && !initialRequest {
			log.Printf("Problems in site - medium.com - at %v", time.Now())
			return articles, nil
		}

		for _, state := range states {
			lenItemsStates := len(state.Data.TagFeed.Items)
			for index, itemState := range state.Data.TagFeed.Items {
				article := model.Article{
					Title:       itemState.Post.Title,
					URL:         itemState.Post.MediumUrl,
					Created:     time.Unix(itemState.Post.FirstPublishedAt/1000, 0),
					Author:      itemState.Post.Creator.Name,
					Description: itemState.Post.ExtendedPreviewContent.Subtitle,
				}
				articles = append(articles, article)
				if index == (lenItemsStates - 1) {
					dateLastState = itemState.Post.FirstPublishedAt / 1000
				}
			}
		}

		initialRequest = false

		// For implement pagination, this required in request
		initNumberState += limitQuantityStates

		if maxDateUnix >= dateLastState {
			break
		}

	}

	return
}

func getStates(url string, initialRequest bool, initNumberState int) (states []StatesMedium, err error) {

	responseBody, err := getBody(initialRequest, initNumberState)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, responseBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Host", "medium.com")
	request.Header.Set("Content-Type", "application/json")

	responce, err := http.DefaultClient.Do(request)
	if err != nil || responce.StatusCode != http.StatusOK {
		return nil, err
	}

	defer responce.Body.Close()

	resBody, err := ioutil.ReadAll(responce.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &states)
	if err != nil {
		return nil, err
	}

	return states, nil

}

func getBody(initialRequest bool, initNumberState int) (*bytes.Buffer, error) {

	var postBody []byte

	paging := fmt.Sprintf(`{"to":"%d","limit":%d}}`, initNumberState, limitQuantityStates)
	if initialRequest {
		paging = fmt.Sprintf(`{"limit":%d}}`, limitQuantityStates)
	}

	textFile, err := os.ReadFile(queryFileName)
	if err != nil {
		return nil, err
	}

	postBody = []byte(fmt.Sprintf(`[{"operationName":"TopicFeedQuery","variables":{"tagSlug":"golang","mode":"NEW",
		"paging": %s,"query":"%s"}]`, paging, string(textFile)))

	return bytes.NewBuffer(postBody), nil

}
