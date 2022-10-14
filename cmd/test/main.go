package main

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/internal/parser/stub"
)

func main() {
	p := stub.NewParser("http://stub.com")

	articles, err := p.ParseAll()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(articles)
}
