package webservice

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/model"
)

var mockNews = []model.Article{
	{
		Title:       "1",
		Created:     time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
		Description: "Short description of the article1",
		URL:         "http://funlink1.com",
		Author:      "m.mikhailov",
	},
}

func TestMessageHandler(t *testing.T) {

	handler := MessageHandler(mockNews)
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/last_news", nil)
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
	expected := `[{"title":"1","created":"2022-01-01T01:01:01.000000001Z","author":"m.mikhailov","summary":"Short description of the article1","url":"http://funlink1.com"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
