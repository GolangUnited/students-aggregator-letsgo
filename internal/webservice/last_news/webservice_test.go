package last_news

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	clog "github.com/indikator/aggregator_lets_go/internal/log/common"

	"github.com/indikator/aggregator_lets_go/model"
)

var id1 = 1
var id2 = 2

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

	handler := ws.GetLastNews(7)

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
	err = json.Unmarshal(rr.Body.Bytes(), resp)
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}

	aGot := db.ConvertFromDbArticles(*resp)
	aExpected := db.ConvertFromDbArticles(expected)

	if !reflect.DeepEqual(aGot, aExpected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			*resp, expected)
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

func TestGetBuildVersion(t *testing.T) {
	err := os.Setenv("TAG", "v0.0.27")
	if err != nil {
		t.Errorf("expected nil got %v", err)
	}
	c := config.NewConfig()
	err = c.SetDataFromFile("../../../tests/configs/webservice/config_buildVersion.yaml")
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

	handler := ws.GetBuildVersion()

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

	responseData, err := ioutil.ReadAll(rr.Body)
	expected := []byte("v0.0.27")

	if !reflect.DeepEqual(responseData, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseData, expected)
	}
}
