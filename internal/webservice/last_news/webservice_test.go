package last_news

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type database struct {
	url string
}

// NewDb create an instance of database
func NewDb(URL string) db.Db {
	return &database{
		url: URL,
	}
}

var id1 = primitive.NewObjectID()
var id2 = primitive.NewObjectID()

func (mockdb *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	return article, nil
}

func (mockdb *database) ReadAllArticles() ([]model.DBArticle, error) {
	articles := []model.DBArticle{
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
	return articles, nil
}

// DBInit creates a new MongoDB client and connect to your running MongoDB server
func (mockdb *database) DBInit() error {

	fmt.Printf("Connected to %s\n", mockdb.url)

	return nil

}

func TestMessageHandler(t *testing.T) {
	ws := NewWebservice("/last_news")
	db := NewDb("test_url")

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
	expected := fmt.Sprintf(`[{"ID":"%s","Title":"test_title","Created":"2022-01-01T02:01:01.000000001Z","Author":"mikhailov.mk","Description":"test article for db","URL":"test_article.com"},{"ID":"%s","Title":"test_title1","Created":"2022-01-01T02:01:01.000000001Z","Author":"mikhailov.mk","Description":"test article1 for db","URL":"test_article1.com"}]`, id1.Hex(), id2.Hex())
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
