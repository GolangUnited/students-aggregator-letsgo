package main

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"github.com/indikator/aggregator_lets_go/model"
)

func main() {
	articles := make([]model.Article, 0)
	webservice.RunServer(articles)
	fmt.Println(articles)
}
