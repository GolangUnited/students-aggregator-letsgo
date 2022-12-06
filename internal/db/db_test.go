package db

import (
	"reflect"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
		ID:          primitive.NewObjectID(),
		Title:       "article 1",
		Created:     dt,
		Description: "Stub article 1",
		URL:         "http://stub.com/1.html",
		Author:      "Stub Stub",
	}

	dbaGot := ConvertToDbArticle(a)

	if dbaGot.Title != dbaExpected.Title {
		t.Errorf("title \"%s\" incorrect, expected \"%s\"", dbaGot.Title, dbaExpected.Title)
	}

	if dbaGot.Created != dbaExpected.Created {
		t.Errorf("created \"%v\" incorrect, expected \"%v\"", dbaGot.Created, dbaExpected.Created)
	}

	if dbaGot.Description != dbaExpected.Description {
		t.Errorf("description \"%s\" incorrect, expected \"%s\"", dbaGot.Description, dbaExpected.Description)
	}

	if dbaGot.URL != dbaExpected.URL {
		t.Errorf("url \"%s\" incorrect, expected \"%s\"", dbaGot.URL, dbaExpected.URL)
	}

	if dbaGot.Author != dbaExpected.Author {
		t.Errorf("author \"%s\" incorrect, expected \"%s\"", dbaGot.Author, dbaExpected.Author)
	}
}

func TestConvertFromDbArticle(t *testing.T) {
	dt := time.Now()

	dba := model.DBArticle{
		ID:          primitive.NewObjectID(),
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

	if !reflect.DeepEqual(*aGot, aExpected) {
		t.Errorf("article incorrect: got %v, expected %v", *aGot, aExpected)
	}
}
