package medium_test

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/parser/medium"
)

func Test_articlesparser_ParseAll(t *testing.T) {

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{name: "Default case", url: "https://medium.com/_/graphql", wantErr: false},
		{name: "Incorect url", url: "https://mediumsssss.com/_/graphql", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			parser := medium.NewParser(tt.url)
			gotArticles, err := parser.ParseAfter(time.Date(2022, time.November, 8, 14, 0, 0, 0, time.UTC))

			if (err != nil) != tt.wantErr {
				t.Errorf("articlesparser.ParseAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotArticles) == 0 && err == nil {
				t.Errorf("articlesparser.ParseAll() responce is empty")
				return
			}

		})
	}
}
