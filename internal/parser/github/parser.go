package github

import (
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/indikator/aggregator_lets_go/internal/log"
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
	a                         = "a"
	href                      = "href"
	parserDatetimeNode        = "relative-time"
	parserDatetimeRef         = "datetime"
	perserAuthorText          = "Auto-Submit: "
	parserDescriptionRef      = "pre"
)

type articlesparser struct {
	url       string
	log       log.Log
	collector *colly.Collector
}

// create an instance of articles parser
func NewParser(cfg parser.Config, l log.Log) parser.ArticlesParser {
	collector := colly.NewCollector()

	if cfg.IsLocal {
		transport := &http.Transport{}
		transport.RegisterProtocol(scheme, http.NewFileTransport(http.Dir(testWebpageLocation)))
		collector.WithTransport(transport)
	}

	return &articlesparser{
		url:       cfg.URL,
		log:       l,
		collector: collector,
	}
}

// parse all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(maxDate time.Time) (articles []model.Article, err error) {
	p.collector.OnHTML(parserArticleContainerRef, func(h *colly.HTMLElement) {
		newArticle, err2 := p.getNewArticle(h)
		if err2 != nil {
			err = err2
			p.log.WriteError(err.Error(), err)
			return
		}
		if !newArticle.Created.After(maxDate) {
			return
		}
		articles = append(articles, newArticle)
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

// get new article description
func (p *articlesparser) getNewArticle(h *colly.HTMLElement) (article model.Article, err error) {
	title, err := p.getTitle(h)
	if err != nil {
		return
	}

	URL, err := p.getAbsoluteURL(h)
	if err != nil {
		return
	}

	description, err := p.getDescription(h)
	if err != nil {
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

	article = model.Article{
		Title:       title,
		URL:         URL,
		Created:     createdAt,
		Author:      author,
		Description: description,
	}
	return
}

// get article title from html element
func (p *articlesparser) getTitle(h *colly.HTMLElement) (title string, err error) {
	if quantity := h.DOM.Find(parserTitleRef).Length(); quantity == 0 {
		err = parser.ErrorArticleTitleNotFound
	} else {
		title = h.DOM.Find(parserTitleRef).Text()
	}
	return
}

// get article absolute url
func (p *articlesparser) getAbsoluteURL(h *colly.HTMLElement) (URL string, err error) {
	if path, ok := h.DOM.Find(a).Attr(href); !ok {
		err = parser.ErrorArticleURLNotFound
	} else {
		URL = h.Request.AbsoluteURL(path)
	}
	return
}

// get article datetime
func (p *articlesparser) getDatetime(h *colly.HTMLElement) (datetime time.Time, err error) {
	if strdate, ok := h.DOM.Find(parserDatetimeNode).Attr(parserDatetimeRef); !ok {
		err = parser.ErrorArticleDatetimeNotFound
	} else {
		datetime, err = time.Parse(dateFormat, strdate)
		if err != nil {
			err = parser.ErrorCannotParseArticleDatetime{OriginError: err}
		}
	}
	return
}

// get article author
func (p *articlesparser) getAuthor(h *colly.HTMLElement) (author string, err error) {
	descriptionText, err := p.getDescription(h)
	if err != nil {
		err = parser.ErrorCannotParseArticleAuthor{OriginError: err}
		return
	}

	descriptionParts := strings.Split(descriptionText, "\n")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], perserAuthorText) {
			idx := strings.Index(descriptionParts[i], perserAuthorText)
			author = descriptionParts[i][idx+len(perserAuthorText):]
			return
		}
	}

	err = parser.ErrorArticleAuthorNotFound
	return
}

// get article description (summary)
func (p *articlesparser) getDescription(h *colly.HTMLElement) (description string, err error) {
	foundDescr := h.DOM.Find(parserDescriptionRef)
	if quantity := foundDescr.Length(); quantity == 0 {
		err = parser.ErrorArticleDescriptionNotFound
	} else {
		description = strings.TrimSpace(foundDescr.Text())
	}
	return
}

func init() {
	parser.RegisterParser("github", NewParser)
}
