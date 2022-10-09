package golang

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	dateFormat = "2006-01-02T15:04:05Z"
)

type articlesparser struct {
	url string
}

// / create an instance of articles parser
func NewParser(URL string) parser.ArticlesParser {
	return &articlesparser{
		url: URL,
	}
}

func (p *articlesparser) getArticleContainerRef() string {
	return "div.js-details-container"
}

func (p *articlesparser) getTitle(h *colly.HTMLElement) string {
	return h.ChildText("a.Link--primary")
}

func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) string {
	return h.Request.AbsoluteURL(h.ChildAttr("a", "href"))
}

func (p *articlesparser) getDatetime(h *colly.HTMLElement) time.Time {
	// receive datetime in format "2022-10-04T17:43:19Z"
	strdate := h.ChildAttr("relative-time", "datetime")
	datetime, _ := time.Parse(dateFormat, strdate)
	return datetime
}

func (p *articlesparser) getDescription(h *colly.HTMLElement) string {
	return h.ChildText("pre.color-fg-muted")
}

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

// / parse all avaibale articles on a web page
func (p *articlesparser) ParseAll() (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML(p.getArticleContainerRef(), func(h *colly.HTMLElement) {
		newArticle := p.getNewArticle(h)
		articles = append(articles, newArticle)
	})

	c.Visit(p.url)

	return
}

// / parse all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML(p.getArticleContainerRef(), func(h *colly.HTMLElement) {
		date := p.getDatetime(h)
		if !date.Before(maxDate) {
			return
		}
		newArticle := p.getNewArticle(h)
		articles = append(articles, newArticle)
	})

	c.Visit(p.url)

	return
}

// / parse n articles with a date less than the given one
func (p *articlesparser) ParseAfterN(maxDate time.Time, n int) (articles []model.Article, err error) {

	c := colly.NewCollector()

	c.OnHTML(p.getArticleContainerRef(), func(h *colly.HTMLElement) {
		date := p.getDatetime(h)
		if !(date.Before(maxDate) && len(articles) < n) {
			return
		}
		newArticle := p.getNewArticle(h)
		articles = append(articles, newArticle)
	})

	c.Visit(p.url)

	return
}

func init() {
	parser.RegisterParser("github", NewParser)
}
