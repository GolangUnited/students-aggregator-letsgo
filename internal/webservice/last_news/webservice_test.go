package last_news

import (
	"encoding/json"
	"fmt"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/config"

	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var id1 = primitive.NewObjectID()
var id2 = primitive.NewObjectID()

func TestGetLastNews(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/webservice/config.yaml")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	ws := NewWebservice()
	err = ws.InitAllByConfig(c)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	// Check the response body is what we expect.
	expected := []model.DBArticle{
		{
			ID:          id1,
			Title:       "test_title",
			Created:     time.Date(2022, 1, 1, 2, 1, 1, 1, time.UTC),
			Author:      "mikhailov.mk",
			Description: "test article for db",
			URL:         "test_article.com",
		},
		{
			ID:          id2,
			Title:       "test_title1",
			Created:     time.Date(2022, 1, 1, 2, 1, 1, 1, time.UTC),
			Author:      "mikhailov.mk",
			Description: "test article1 for db",
			URL:         "test_article1.com",
		},
	}
	newDb := ws.Db()
	for _, article := range expected {
		newDb.WriteArticle(&article)
	}

	handler, err := ws.GetLastNews(7)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", c.WebService.Handle, nil)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := &[]model.DBArticle{}
	err = json.Unmarshal([]byte(rr.Body.String()), resp)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	if !reflect.DeepEqual(*resp, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp, expected)
	}

	_, err = ws.GetLastNews(-1)
	if !reflect.DeepEqual(err, fmt.Errorf("invalid number of days %d", -1)) {
		t.Errorf("expected %v got %v", fmt.Errorf("invalid number of days %d", -1), err)
	}

}

func TestPort(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/webservice/config.yaml")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	ws := NewWebservice()
	err = ws.InitAllByConfig(c)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	port := ws.Port()
	if port != 8080 {
		t.Errorf("unexpected port")
	}
}

func TestLogger(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/webservice/config.yaml")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	ws := NewWebservice()
	err = ws.InitAllByConfig(c)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	logger := ws.Logger()
	l, err := clog.GetLog(c.WebService.Log)
	if !reflect.DeepEqual(logger, l) {
		t.Errorf("unexpected logger")
	}
}

func TestInitAllByConfig(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/webservice/config_getlog.yaml")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	ws := NewWebservice()
	err = ws.InitAllByConfig(c)
	if err.Error() != "unknown log type undefined" {
		t.Errorf("expected %v got %v", "unknown log type undefined", err.Error())
	}

	c = config.NewConfig()
	err = c.SetDataFromFile("../../../tests/configs/webservice/config_getdb.yaml")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	err = ws.InitAllByConfig(c)
	if err.Error() != "unknown dbms undefined" {
		t.Errorf("expected %v got %v", "unknown dbms undefined", err.Error())
	}

}
