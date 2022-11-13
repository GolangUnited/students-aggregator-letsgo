package github

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

// config consts
const (
	dateFormat          = "2006-01-02T15:04:05Z"
	scheme              = "file"
	testWebpageLocation = "../../../tests/data/parser/github/"
)

// parser consts
const (
	parserArticleContainerRef = "div.js-details-container"
	parserTitleRef            = "a.Link--primary"
	parserAbsoluteUrlNode     = "a"
	parserAbsoluteUrlRef      = "href"
	parserDatetimeNode        = "relative-time"
	parserDatetimeRef         = "datetime"
	perserAuthorText          = "Auto-Submit: "
	parserDescriptionRef      = "pre.color-fg-muted"
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
	p.collector.OnHTML(parserArticleContainerRef, func(h *colly.HTMLElement) {
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
		Author:      p.getAuthor(h),
		Description: p.getDescription(h),
	}
	return newArticle
}

// get article title from html element
func (p *articlesparser) getTitle(h *colly.HTMLElement) string {
	return h.ChildText(parserTitleRef)
}

// get article absolute url
func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) string {
	return h.Request.AbsoluteURL(h.ChildAttr(parserAbsoluteUrlNode, parserAbsoluteUrlRef))
}

// get article datetime
func (p *articlesparser) getDatetime(h *colly.HTMLElement) time.Time {
	// receive datetime in format "2022-10-04T17:43:19Z"
	strdate := h.ChildAttr(parserDatetimeNode, parserDatetimeRef)
	datetime, _ := time.Parse(dateFormat, strdate)
	return datetime
}

// get article author
func (p *articlesparser) getAuthor(h *colly.HTMLElement) string {
	descriptionText := h.ChildText(parserDescriptionRef)
	descriptionParts := strings.Split(descriptionText, "\n")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], perserAuthorText) {
			idx := strings.Index(descriptionParts[i], perserAuthorText)
			return descriptionParts[i][idx+len(perserAuthorText):]
		}
	}

	return ""
}

// get article description (summary)
func (p *articlesparser) getDescription(h *colly.HTMLElement) string {
	return h.ChildText(parserDescriptionRef)
}

func init() {
	parser.RegisterParser("github", NewParser)
}
