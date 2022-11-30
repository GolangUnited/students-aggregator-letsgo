package godev

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	dateFormat          = "2 January 2006"
	scheme              = "file"
	testWebPageLocation = "../../../tests/data/parser/godev/"

	parserName = "go.dev"

	articleContainerTag   = "p.blogtitle"
	articleTitleTag       = "a[href]"
	articleDatetimeTag    = "span.date"
	articleAuthorTag      = "span.author"
	articleDescriptionTag = "p.blogsummary"
	a                     = "a"
	href                  = "href"
	allArticlesHref       = "file://./blog/all"
)

type articlesparser struct {
	url       string
	log       log.Log
	collector *colly.Collector
}

// NewParser creates an instance of articles parser
func NewParser(cfg parser.Config, lg log.Log) parser.ArticlesParser {

	collector := colly.NewCollector()

	if cfg.IsLocal {
		transport := &http.Transport{}
		transport.RegisterProtocol(scheme, http.NewFileTransport(http.Dir(testWebPageLocation)))
		collector.WithTransport(transport)
	}

	return &articlesparser{
		url:       cfg.URL,
		log:       lg,
		collector: collector,
	}
}

// ParseAfter parses all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(date time.Time) (articles []model.Article, err error) {

	articleContainerRef := p.getArticleContainerRef()

	p.collector.OnHTML(articleContainerRef, func(h *colly.HTMLElement) {
		article, er := p.getNewArticle(h)
		if er != nil {
			err = er
			p.log.WriteError(err.Error(), err)
			return
		}
		if !article.Created.After(date) {
			return
		}
		articles = append(articles, article)
	})

	p.collector.OnError(func(r *colly.Response, er error) {
		if strings.TrimSpace(er.Error()) == parser.ErrorMessage {
			err = parser.ErrorWebPageCannotBeDelivered{URL: r.Request.URL.String(), StatusCode: r.StatusCode}
		} else {
			err = parser.ErrorUnknown{OriginError: er}
		}
		p.log.WriteError(err.Error(), err)
	})

	p.collector.Visit(p.url)

	return
}

// getArticleContainerRef gets parseable articles html container
func (p *articlesparser) getArticleContainerRef() string {
	return articleContainerTag
}

// getTitle gets an article title from html element
func (p *articlesparser) getTitle(h *colly.HTMLElement) (title string, err error) {

	if quantity := h.DOM.Find(articleTitleTag).Length(); quantity == 0 {
		err = parser.ErrorArticleTitleNotFound
	} else {
		title = h.DOM.Find(articleTitleTag).Text()
	}

	return
}

// getAbsoluteURL gets an article absolute url
func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) (URL string, err error) {

	if path, ok := h.DOM.Find(a).Attr(href); !ok {
		err = parser.ErrorArticleURLNotFound
	} else {
		URL = h.Request.AbsoluteURL(path)
	}

	return
}

// getDatetime gets an article datetime
func (p *articlesparser) getDatetime(h *colly.HTMLElement) (datetime time.Time, err error) {

	if quantity := h.DOM.Find(articleDatetimeTag).Length(); quantity == 0 {
		err = parser.ErrorArticleDatetimeNotFound
	} else {
		strdate := h.DOM.Find(articleDatetimeTag).Text()
		datetime, err = time.Parse(dateFormat, strdate)
		if err != nil {
			err = parser.ErrorCannotParseArticleDatetime{OriginError: err}
		}
	}

	return
}

// getAuthor gets an article author
func (p *articlesparser) getAuthor(h *colly.HTMLElement) (author string, err error) {

	if quantity := h.DOM.Find(articleAuthorTag).Length(); quantity == 0 {
		err = parser.ErrorArticleAuthorNotFound
	} else {
		author = h.DOM.Find(articleAuthorTag).Text()
	}

	return
}

// getDescription gets an article description (summary)
func (p *articlesparser) getDescription(h *colly.HTMLElement) (description string, err error) {

	if quantity := h.DOM.NextFiltered(articleDescriptionTag).Length(); quantity == 0 {
		err = parser.ErrorArticleDescriptionNotFound
	} else {
		description = strings.TrimSpace(h.DOM.NextFiltered(articleDescriptionTag).Text())
	}

	return
}

// getNewArticle gets a new article
func (p *articlesparser) getNewArticle(h *colly.HTMLElement) (article model.Article, err error) {

	title, err := p.getTitle(h)
	if err != nil {
		return
	}

	URL, err := p.getAbsoluteURL(h)
	if err != nil || URL == allArticlesHref {
		return
	}

	author, err := p.getAuthor(h)
	if err != nil {
		return
	}

	createdAt, err := p.getDatetime(h)
	if err != nil {
		return
	}

	description, err := p.getDescription(h)
	if err != nil {
		return
	}

	article = model.Article{
		Title:       title,
		URL:         URL,
		Created:     createdAt,
		Author:      author,
		Description: description,
	}

	return
}

func init() {
	parser.RegisterParser(parserName, NewParser)
}
