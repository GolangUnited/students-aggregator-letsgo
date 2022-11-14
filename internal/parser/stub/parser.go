package stub

import (
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/model"
)

type articlesparser struct {
	url string
}

func prepare(url string, date time.Time) []model.Article {
	var articles = []model.Article{{
		Title:       "article 1",
		Created:     date,
		Description: "Stub article 1",
		URL:         url + "/1.html",
		Author:      "Stub Stub",
	}, {
		Title:       "article 2",
		Created:     date.Add(-24 * time.Hour),
		Description: "Stub article 2",
		URL:         url + "/2.html",
		Author:      "Stub Stub",
	}, {
		Title:       "article 3",
		Created:     date.Add(-48 * time.Hour),
		Description: "Stub article 3",
		URL:         url + "/3.html",
		Author:      "Stub Stub",
	},
	}

	return articles
}

// create an instance of articles parser
func NewParser(cfg parser.Config) parser.ArticlesParser {
	return &articlesparser{
		url: cfg.URL,
	}
}

// parse all articles that were created earler than the target date
func (p *articlesparser) ParseAfter(maxDate time.Time) ([]model.Article, error) {
	articles := prepare(p.url, maxDate)

	articles2 := make([]model.Article, 0)

	for _, a := range articles {
		if a.Created.Before(maxDate) {
			articles2 = append(articles2, a)
		}
	}

	return articles2, nil
}

func init() {
	parser.RegisterParser("stub", NewParser)
}
