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
	hostService   = "medium.com"
)

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

type ArticlesParser struct {
	Client        *http.Client
	Url           string
	Host          string
	QueryFileName string
}

// create an instance of articles parser
func NewParser(cfg parser.Config) parser.ArticlesParser {

	return &ArticlesParser{
		Client:        &http.Client{},
		Url:           cfg.URL,
		Host:          hostService,
		QueryFileName: queryFileName,
	}

}

func init() {
	parser.RegisterParser(hostService, NewParser)
}

func (p *ArticlesParser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	var (
		initNumberState int
		firstCreated    time.Time
	)

	initialRequest := true

	for true {

		states, err := getStates(p, initialRequest, initNumberState)
		if err != nil {
			return nil, err
		}

		if len(states[0].Data.TagFeed.Items) == 0 && !initialRequest {
			log.Printf("Problems in site - %s", p.Host)
			return articles, nil
		}

		for _, state := range states {
			for _, itemState := range state.Data.TagFeed.Items {

				created := time.Unix(itemState.Post.FirstPublishedAt/1000, 0)
				if firstCreated.IsZero() {
					firstCreated = created
				}
				if created.After(maxDate) {
					article := model.Article{
						Title:       itemState.Post.Title,
						URL:         itemState.Post.MediumUrl,
						Created:     created,
						Author:      itemState.Post.Creator.Name,
						Description: itemState.Post.ExtendedPreviewContent.Subtitle,
					}
					articles = append(articles, article)
				}

				if (firstCreated.Sub(created).Hours()/24 >= 7) || len(articles) >= 100 {
					return articles, nil
				}
			}
		}

		if initialRequest {
			initialRequest = false
		}

		// For implement pagination, this required in request
		initNumberState += limitQuantityStates

	}

	return
}

func getStates(p *ArticlesParser, initialRequest bool, initNumberState int) (states []StatesMedium, err error) {

	requestBody, err := getBodyRequest(p, initialRequest, initNumberState)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, p.Url, requestBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Host", p.Host)
	request.Header.Set("Content-Type", "application/json")

	response, err := p.Client.Do(request)
	if err != nil || response.StatusCode != http.StatusOK {
		return nil, err
	}

	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &states)
	if err != nil {
		return nil, err
	}

	return states, nil

}

func getBodyRequest(p *ArticlesParser, initialRequest bool, initNumberState int) (*bytes.Buffer, error) {

	var postBody []byte

	paging := fmt.Sprintf(`{"to":"%d","limit":%d}}`, initNumberState, limitQuantityStates)
	if initialRequest {
		paging = fmt.Sprintf(`{"limit":%d}}`, limitQuantityStates)
	}

	textFile, err := os.ReadFile(p.QueryFileName)
	if err != nil {
		return nil, err
	}

	postBody = []byte(fmt.Sprintf(`[{"operationName":"TopicFeedQuery","variables":{"tagSlug":"golang","mode":"NEW",
		"paging": %s,"query":"%s`, paging, string(textFile)))

	return bytes.NewBuffer(postBody), nil

}
