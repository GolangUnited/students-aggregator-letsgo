package godev

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	dateFormat = "2 January 2006"
)

type articlesparser struct {
	url string
}

/// create an instance of articles parser
func NewParser(URL string) parser.ArticlesParser {
	return &articlesparser{
		url: URL,
	}
}

/// parse all avaibale articles on a web page
func (p *articlesparser) ParseAll() (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		article := model.Article{
			Title:       h.ChildText("a[href]"),
			URL:         h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Created:     date,
			Author:      h.ChildText("span.author"),
			Description: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		articles = append(articles, article)
	})

	c.Visit(p.url)

	return
}

/// parse all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		if !date.Before(maxDate) {
			return
		}
		article := model.Article{
			Title:       h.ChildText("a[href]"),
			URL:         h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Created:     date,
			Author:      h.ChildText("span.author"),
			Description: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		articles = append(articles, article)
	})

	c.Visit(p.url)

	return
}

/// parse n articles with a date less than the given one
func (p *articlesparser) ParseAfterN(maxDate time.Time, n int) (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		if !(date.Before(maxDate) && len(articles) < n) {
			return
		}
		article := model.Article{
			Title:       h.ChildText("a[href]"),
			URL:         h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Created:     date,
			Author:      h.ChildText("span.author"),
			Description: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		articles = append(articles, article)
	})

	c.Visit(p.url)

	return
}

func init() {
	parser.RegisterParser("go.dev", NewParser)
}
