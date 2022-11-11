package aggregator

import (
	"reflect"
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
)

const (
	configFilePath = "../../tests/configs/aggregator/config.yaml"
)

func TestWorkWithStubParserAndDb(t *testing.T) {

	c := config.NewConfig()

	err := c.SetDataFromFile(configFilePath)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	a := NewAggregator()

	err = a.InitAllByConfig(c)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = a.Execute()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	articlesFromParser, err := a.parsers[0].ParseAfter(a.lastCheckDatetime)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	articles, err := a.db.ReadAllArticles()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	var articlesFromDb []model.Article

	for _, article := range articles {
		dbArticle := db.ConvertFromDbArticle(article)

		articlesFromDb = append(articlesFromDb, *dbArticle)
	}

	if !reflect.DeepEqual(articlesFromDb, articlesFromParser) {
		t.Log(articlesFromDb, articlesFromParser)

		t.Errorf("articles received from a database are different from articles received from a parser")
	}
}
