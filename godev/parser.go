package godev

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const (
	filePath   = "posts.json"
	dateFormat = "2 January 2006"
)

type NewsParser interface {
	ParseAll() (posts []Post, err error)
}

type Post struct {
	Title   string    `json:"title"`
	Href    string    `json:"href"`
	Date    time.Time `json:"date"`
	Author  string    `json:"author"`
	Summary string    `json:"summary"`
}

type parser struct {
	url string
}

func NewParser(URL string) NewsParser {
	return &parser{
		url: URL,
	}
}

func (p *parser) ParseAll() (posts []Post, err error) {

	c := colly.NewCollector()

	c.OnHTML("p.blogtitle", func(h *colly.HTMLElement) {
		date, _ := time.Parse(dateFormat, h.ChildText("span.date"))
		post := Post{
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
