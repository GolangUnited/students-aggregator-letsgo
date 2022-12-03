package medium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
	"github.com/pkg/errors"
)

const (
	limitQuantityStates = 25
	// The file contains a request to the service database
	queryFileName = "./etc/queryInMedium"
	hostService   = "medium.com"
	// The request returns a maximum of 25 articles at a time
	// and we do not need more than a hundred articles in one call to the parser
	iterationCount = 4
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
	Client         *http.Client
	Url            string
	Host           string
	QueryFileName  string
	IterationCount int
	Log            log.Log
	LocalLaunch    bool
}

// create an instance of articles parser
func NewParser(cfg parser.Config, l log.Log) parser.ArticlesParser {

	return &ArticlesParser{
		Client:         &http.Client{},
		Url:            cfg.URL,
		Host:           hostService,
		QueryFileName:  queryFileName,
		IterationCount: iterationCount,
		Log:            l,
		LocalLaunch:    false,
	}

}

func init() {
	parser.RegisterParser(hostService, NewParser)
}

func (p *ArticlesParser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	var (
		initNumberState int
		states          []StatesMedium
		article         model.Article
	)

	initialRequest := true

	for i := 0; i < p.IterationCount; i++ {

		states, err = p.getStates(initialRequest, initNumberState)
		if err != nil {
			return nil, err
		}

		if len(states[0].Data.TagFeed.Items) == 0 {
			err = errors.Errorf("Problems in site - %s", p.Host)
			err = parser.ErrorUnknown{OriginError: err}
			p.Log.WriteError(err.Error(), err)
			if !p.LocalLaunch {
				err = nil
			}
			return articles, err
		}

		for _, state := range states {
			for _, itemState := range state.Data.TagFeed.Items {

				article, err = p.getNewArticle(&itemState)
				if err != nil {
					p.Log.WriteError(err.Error(), err)
					if !p.LocalLaunch {
						err = nil
					}
					continue
				}

				if !article.Created.After(maxDate) {
					continue
				}
				articles = append(articles, article)

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

func (p *ArticlesParser) getStates(initialRequest bool, initNumberState int) (states []StatesMedium, err error) {

	requestBody, err := p.getBodyRequest(initialRequest, initNumberState)
	if err != nil {
		err = parser.ErrorUnknown{OriginError: err}
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, p.Url, requestBody)
	if err != nil {
		err = parser.ErrorUnknown{OriginError: err}
		return nil, err
	}
	request.Header.Set("Host", p.Host)
	request.Header.Set("Content-Type", "application/json")

	response, err := p.Client.Do(request)
	if err != nil {
		err = parser.ErrorUnknown{OriginError: err}
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			err = parser.ErrorWebPageCannotBeDelivered{URL: request.URL.String(), StatusCode: response.StatusCode}
		} else {
			err = parser.ErrorUnknown{OriginError: err}
		}
		return nil, err
	}

	defer response.Body.Close()

	resBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = parser.ErrorUnknown{OriginError: err}
		return nil, err
	}

	err = json.Unmarshal(resBody, &states)
	if err != nil {
		err = parser.ErrorUnknown{OriginError: err}
		return nil, err
	}

	return states, nil

}

func (p *ArticlesParser) getBodyRequest(initialRequest bool, initNumberState int) (*bytes.Buffer, error) {

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

// getNewArticle gets a new article
func (p *ArticlesParser) getNewArticle(itemState *Item) (article model.Article, err error) {

	title := itemState.Post.Title
	if len(title) == 0 {
		err = parser.ErrorArticleTitleNotFound
		return
	}

	url := itemState.Post.MediumUrl
	if len(url) == 0 {
		err = parser.ErrorArticleURLNotFound
		return
	}

	author := itemState.Post.Creator.Name
	if len(author) == 0 {
		err = parser.ErrorArticleAuthorNotFound
		return
	}

	createdAt := time.Unix(itemState.Post.FirstPublishedAt/1000, 0)
	if itemState.Post.FirstPublishedAt == 0 {
		err = parser.ErrorArticleDatetimeNotFound
		return
	}

	description := itemState.Post.ExtendedPreviewContent.Subtitle
	if len(description) == 0 {
		err = parser.ErrorArticleDescriptionNotFound
		return
	}

	article = model.Article{
		Title:       title,
		URL:         url,
		Created:     createdAt,
		Author:      author,
		Description: description,
	}

	return
}
