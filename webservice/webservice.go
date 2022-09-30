package webservice

import (
	"fmt"
	"net/http"
)

func MessageHandler(news []DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, article := range news {
			fmt.Fprintf(w, "ID : %d\nDatetime : %v\nDescription : %s\nLink : %s\nPage file : %s\n", article.Id, article.Datetime, article.Description, article.Link, article.PageHtml)
		}
	})
}
