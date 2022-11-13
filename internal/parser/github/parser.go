package github

import (
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	dateFormat          = "2006-01-02T15:04:05Z"
	scheme              = "file"
	testWebpageLocation = "../../../tests/data/parser/github/"
)

type articlesparser struct {
	url       string
	collector *colly.Collector
}

// create an instance of articles parser
func NewParser(cfg parser.Config) parser.ArticlesParser {
	collector := colly.NewCollector()

	if cfg.IsLocal {
		transport := &http.Transport{}
		transport.RegisterProtocol(scheme, http.NewFileTransport(http.Dir(testWebpageLocation)))
		collector.WithTransport(transport)
	}

	return &articlesparser{
		url:       cfg.URL,
		collector: collector,
	}
}

// parse all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {
	articleContainerRef := p.getArticleContainerRef()
	p.collector.OnHTML(articleContainerRef, func(h *colly.HTMLElement) {
		date := p.getDatetime(h)
		if !date.After(maxDate) {
			return
		}
		newArticle := p.getNewArticle(h)
		articles = append(articles, newArticle)
	})

	p.collector.Visit(p.url)

	return
}

// get new article description
func (p *articlesparser) getNewArticle(h *colly.HTMLElement) model.Article {
	newArticle := model.Article{
		Title:       p.getTitle(h),
		URL:         p.getAbsoluteURL(h),
		Created:     p.getDatetime(h),
		Author:      "",
		Description: p.getDescription(h),
	}
	return newArticle
}

// get parseable articles html container
func (p *articlesparser) getArticleContainerRef() string {
	return "div.js-details-container"
}

// get article title from html element
func (p *articlesparser) getTitle(h *colly.HTMLElement) string {
	return h.ChildText("a.Link--primary")
}

// get article absolute url
func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) string {
	return h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
}

// get article datetime
func (p *articlesparser) getDatetime(h *colly.HTMLElement) time.Time {
	// receive datetime in format "2022-10-04T17:43:19Z"
	strdate := h.ChildAttr("relative-time", "datetime")
	datetime, _ := time.Parse(dateFormat, strdate)
	return datetime
}

// get article description (summary)
func (p *articlesparser) getDescription(h *colly.HTMLElement) string {
	return h.ChildText("pre.color-fg-muted")
}

func init() {
	parser.RegisterParser("github", NewParser)
}
