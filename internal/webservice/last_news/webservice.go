package last_news

import (
	"encoding/json"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"log"
	"net/http"
)

type webService struct {
	handles map[string]webservice.Config
}

func NewWebservice(config map[string]webservice.Config) webservice.Webservice {
	return &webService{
		handles: config,
	}
}

// GetLastNews godoc
// @Summary Retrieves last news from last 7 days
// @Description Get array of Articles for the last 7 days
// @Tags news
// @Produce json
// @Success 200 {object} []model.DBArticle
// @Router /last_news [get]
func (ws *webService) GetLastNews(db db.Db, nDays int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		news, err := db.ReadArticles(nDays)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		newsJson, err := json.Marshal(news)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(newsJson)
	})
}
