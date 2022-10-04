package main

import (
	"aggregator/parsers/stub"
	"fmt"
	"log"
)

func main() {

	p := stub.NewStubParser()

	err := p.Init(map[string]string{"Url": "123"})

	if err != nil {
		log.Fatal(err)
	}

	articles, err := p.ParseAll()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(articles)
}
