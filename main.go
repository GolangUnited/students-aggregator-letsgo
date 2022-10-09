package main

import (
	"fmt"
	"log"

	"github.com/indikator/aggregator_lets_go/config"
)

func main() {
	c := config.NewConfig()

	err := c.Read("aggregator.yaml")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
