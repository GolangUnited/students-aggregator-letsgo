package main

import (
	"fmt"
	"log"

	//"aggregator_lets_go/config"
	"aggregator/config"
)

func main() {
	c := config.NewConfig()

	err := c.Read("aggregator.yaml")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)
}
