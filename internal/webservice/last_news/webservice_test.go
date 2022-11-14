package last_news

import (
	"encoding/json"
	"github.com/indikator/aggregator_lets_go/internal/db/stub"
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

func TestMessageHandler(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../tests/configs/webservice/config.yaml")
	if err != nil {
		return
	}
	err = c.Read()
	if err != nil {
		return
	}
	ws := NewWebservice(c.WebService)
	newDb := stub.NewDb(c.Database)

	handler := ws.MessageHandler(newDb)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", c.WebService.Handle, nil)
	if err != nil {
		t.Fatal(err)
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

	resp := &[]model.DBArticle{}
	err = json.Unmarshal([]byte(rr.Body.String()), resp)
	if err != nil {
		return
	}

	if !reflect.DeepEqual(resp, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			resp, expected)
	}
}
