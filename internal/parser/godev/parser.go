package godev

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	dateFormat      = "2 January 2006"
	scheme          = "file"
	webpageLocation = "../../../tests/data/parser/godev/"

	parserName = "go.dev"

	articleContainerTag   = "p.blogtitle"
	articleTitleTag       = "a[href]"
	articleDatetimeTag    = "span.date"
	articleAuthorTag      = "span.author"
	articleDescriptionTag = "p.blogsummary"
	a                     = "a"
	href                  = "href"
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
		transport.RegisterProtocol(scheme, http.NewFileTransport(http.Dir(webpageLocation)))
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
		article := p.getNewArticle(h)
		if !article.Created.After(maxDate) {
			return
		}
		articles = append(articles, article)
	})

	p.collector.Visit(p.url)

	return
}

// get parseable articles html container
func (p *articlesparser) getArticleContainerRef() string {
	return articleContainerTag
}

// get an article title from html element
func (p *articlesparser) getTitle(h *colly.HTMLElement) string {
	return h.ChildText(articleTitleTag)
}

// get an article absolute url
func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) string {
	return h.Request.AbsoluteURL(h.ChildAttr(a, href))
}

// get an article datetime
func (p *articlesparser) getDatetime(h *colly.HTMLElement) time.Time {
	strdate := h.ChildText(articleDatetimeTag)
	datetime, _ := time.Parse(dateFormat, strdate)
	return datetime
}

// get an article author
func (p *articlesparser) getAuthor(h *colly.HTMLElement) string {
	return h.ChildText(articleAuthorTag)
}

// get an article description (summary)
func (p *articlesparser) getDescription(h *colly.HTMLElement) string {
	return strings.TrimSpace(h.DOM.NextFiltered(articleDescriptionTag).Text())
}

// get a new article
func (p *articlesparser) getNewArticle(h *colly.HTMLElement) model.Article {
	newArticle := model.Article{
		Title:       p.getTitle(h),
		URL:         p.getAbsoluteURL(h),
		Created:     p.getDatetime(h),
		Author:      p.getAuthor(h),
		Description: p.getDescription(h),
	}
	return newArticle
}

func init() {
	parser.RegisterParser(parserName, NewParser)
}
