package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

func checkArticlesEquals(t *testing.T, aGot model.Article, aExpected model.Article) {
	if !reflect.DeepEqual(aGot, aExpected) {
		t.Errorf("article incorrect: got %v, expected %v", aGot, aExpected)
	}
}

func checkDBArticlesEquals(t *testing.T, aGot model.DBArticle, aExpected model.DBArticle) {
	if aGot.Title != aExpected.Title {
		t.Errorf("title \"%s\" incorrect, expected \"%s\"", aGot.Title, aExpected.Title)
	}

	if aGot.Created != aExpected.Created {
		t.Errorf("created \"%v\" incorrect, expected \"%v\"", aGot.Created, aExpected.Created)
	}

	if aGot.Description != aExpected.Description {
		t.Errorf("description \"%s\" incorrect, expected \"%s\"", aGot.Description, aExpected.Description)
	}

	if aGot.URL != aExpected.URL {
		t.Errorf("url \"%s\" incorrect, expected \"%s\"", aGot.URL, aExpected.URL)
	}

	if aGot.Author != aExpected.Author {
		t.Errorf("author \"%s\" incorrect, expected \"%s\"", aGot.Author, aExpected.Author)
	}
}

func TestConvertToDbArticle(t *testing.T) {
	dt := time.Now()
	a := model.Article{
		Title:       "article 1",
		Created:     dt,
		Description: "Stub article 1",
		URL:         "http://stub.com/1.html",
		Author:      "Stub Stub",
	}

	dbaExpected := model.DBArticle{
		ID:          1,
		Title:       "article 1",
		Created:     dt,
		Description: "Stub article 1",
		URL:         "http://stub.com/1.html",
		Author:      "Stub Stub",
	}

	dbaGot := ConvertToDbArticle(a)

	checkDBArticlesEquals(t, *dbaGot, dbaExpected)
}

func TestConvertFromDbArticle(t *testing.T) {
	dt := time.Now()

	dba := model.DBArticle{
		ID:          1,
		Title:       "article 1",
		Created:     dt,
		Description: "Stub article 1",
		URL:         "http://stub.com/1.html",
		Author:      "Stub Stub",
	}

	aExpected := model.Article{
		Title:       "article 1",
		Created:     dt,
		Description: "Stub article 1",
		URL:         "http://stub.com/1.html",
		Author:      "Stub Stub",
	}

	aGot := ConvertFromDbArticle(dba)

	checkArticlesEquals(t, *aGot, aExpected)
}

func TestConvertToDbArticles(t *testing.T) {
	dt := time.Now()
	a := []model.Article{
		{
			Title:       "article 1",
			Created:     dt,
			Description: "Stub article 1",
			URL:         "http://stub.com/1.html",
			Author:      "Stub Stub",
		},
		{
			Title:       "article 2",
			Created:     dt.Add(time.Second),
			Description: "Stub article 2",
			URL:         "http://stub.com/2.html",
			Author:      "Stub Stub 2",
		},
	}

	dbaExpected := []model.DBArticle{
		{
			ID:          1,
			Title:       "article 1",
			Created:     dt,
			Description: "Stub article 1",
			URL:         "http://stub.com/1.html",
			Author:      "Stub Stub",
		},
		{
			ID:          2,
			Title:       "article 2",
			Created:     dt.Add(time.Second),
			Description: "Stub article 2",
			URL:         "http://stub.com/2.html",
			Author:      "Stub Stub 2",
		},
	}

	dbaGot := ConvertToDbArticles(a)

	if len(dbaGot) != len(dbaExpected) {
		t.Errorf("length %d incorrect, expected %d", len(dbaGot), len(dbaExpected))
	}

	for i := 0; i < len(dbaGot); i++ {
		checkDBArticlesEquals(t, dbaGot[i], dbaExpected[i])
	}
}

func TestConvertFromDbArticles(t *testing.T) {
	dt := time.Now()

	dba := []model.DBArticle{
		{
			ID:          1,
			Title:       "article 1",
			Created:     dt,
			Description: "Stub article 1",
			URL:         "http://stub.com/1.html",
			Author:      "Stub Stub",
		},
		{
			ID:          2,
			Title:       "article 2",
			Created:     dt.Add(time.Second),
			Description: "Stub article 2",
			URL:         "http://stub.com/2.html",
			Author:      "Stub Stub 2",
		},
	}

	aExpected := []model.Article{
		{
			Title:       "article 1",
			Created:     dt,
			Description: "Stub article 1",
			URL:         "http://stub.com/1.html",
			Author:      "Stub Stub",
		},
		{
			Title:       "article 2",
			Created:     dt.Add(time.Second),
			Description: "Stub article 2",
			URL:         "http://stub.com/2.html",
			Author:      "Stub Stub 2",
		},
	}

	aGot := ConvertFromDbArticles(dba)

	if len(aGot) != len(aExpected) {
		t.Errorf("length %d incorrect, expected %d", len(aGot), len(aExpected))
	}

	for i := 0; i < len(aGot); i++ {
		checkArticlesEquals(t, aGot[i], aExpected[i])
	}
}
