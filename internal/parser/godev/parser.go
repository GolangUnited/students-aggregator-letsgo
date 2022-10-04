package godev

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	filePath   = "posts.json"
	dateFormat = "2 January 2006"
)

type newsparser struct {
	url string
}

/// create an instance of the parser
func NewParser(URL string) parser.NewsParser {
	return &newsparser{
		url: URL,
	}
}

/// parse all avaibale posts on a web page
func (p *newsparser) ParseAll() (posts []model.Post, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		post := model.Post{
			Title:   h.ChildText("a[href]"),
			Href:    h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Date:    date,
			Author:  h.ChildText("span.author"),
			Summary: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		posts = append(posts, post)
	})

	c.Visit(p.url)

	return
}

/// parse posts with a date less than the given one
func (p *newsparser) ParseBefore(maxDate time.Time) (posts []model.Post, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		if !date.Before(maxDate) {
			return
		}
		post := model.Post{
			Title:   h.ChildText("a[href]"),
			Href:    h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Date:    date,
			Author:  h.ChildText("span.author"),
			Summary: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		posts = append(posts, post)
	})

	c.Visit(p.url)

	return
}

/// parse n posts with a date less than the given one
func (p *newsparser) ParseBeforeN(maxDate time.Time, n int) (posts []model.Post, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		if !(date.Before(maxDate) && len(posts) < n) {
			return
		}
		post := model.Post{
			Title:   h.ChildText("a[href]"),
			Href:    h.Request.AbsoluteURL(h.ChildAttr("a", "href")),
			Date:    date,
			Author:  h.ChildText("span.author"),
			Summary: strings.TrimSpace(h.DOM.NextFiltered("p.blogsummary").Text()),
		}
		posts = append(posts, post)
	})

	c.Visit(p.url)

	return
}
