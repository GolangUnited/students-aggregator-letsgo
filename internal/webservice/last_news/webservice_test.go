package last_news

import (
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/db/mongo"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMessageHandler(t *testing.T) {
	c := config.NewConfig()
	err := c.Read("../../../etc/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	ws := NewWebservice("/last_news")
	db := mongo.NewDb(c.Database.Url)

	handler := ws.MessageHandler(db)
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
	expected := `[{"ID":"634af235bc6ca8ab6a2af4d3","Title":"test_title","Created":"2022-01-01T01:01:01Z","Author":"mikhailov.mk","Description":"test article for db","URL":"test_article.com"},{"ID":"634af2803e95239b839a93f0","Title":"test_title1","Created":"2022-01-01T01:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"},{"ID":"634af2803e95239b839a93f1","Title":"test_title2","Created":"2022-01-01T01:01:01Z","Author":"mikhailov.mk","Description":"test article2 for db","URL":"test_article2.com"},{"ID":"634afb50ec373cdf2e238503","Title":"test lerik1","Created":"2022-01-01T02:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"},{"ID":"634afc2c8e81c3431f0591b3","Title":"test lerik1","Created":"2022-01-01T02:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"},{"ID":"634b11715e4c2a42650f922f","Title":"test lerik1","Created":"2022-01-01T02:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"},{"ID":"634b117c85ef84f5420f26a2","Title":"test lerik1","Created":"2022-01-01T02:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"},{"ID":"634b11ddb1964d53dc1f53c5","Title":"test lerik1","Created":"2022-01-01T02:01:01Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
