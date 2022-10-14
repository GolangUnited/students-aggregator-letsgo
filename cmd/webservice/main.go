package main

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"github.com/indikator/aggregator_lets_go/model"
)

func main() {
	arcticles := make([]model.Article, 0)
	webservice.RunServer(arcticles)
	fmt.Println(arcticles)
}
