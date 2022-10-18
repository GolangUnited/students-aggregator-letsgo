package main

import (
	"fmt"
	"log"
	"os"

	//"time"

	"github.com/indikator/aggregator_lets_go/internal/config"
	stub "github.com/indikator/aggregator_lets_go/internal/parser/stub"
)

func main() {
	p := stub.NewParser("http://stub.com")
	fmt.Println(p)

	file, er := os.ReadFile("..\\aggregator\\config.yaml")
	if er != nil {
		log.Fatal(er)
	}
	fmt.Println(file)
	c := config.NewConfig()

	s := `# Project Aggregator YAML
aggregator:
  nothing:

database:
  url: mongodb://localhost:27018/

webservice:
  port: 8080

parsers:
- github:
    url: https://github.com/golang/go/tags
- go.dev:
    url: https://go.dev/blog
- medium.com:
    url: https://medium.com/_/graphql`
	fmt.Println([]byte(s))
	err := c.SetData([]byte(s))

	if err != nil {
		log.Fatal(err)
	}

	err = c.Read()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(c)

	//articles, err := p.ParseAll()
	//articles, err := p.ParseAfter(time.Now().Add(-time.Hour))
	// articles, err := p.ParseAfterN(time.Now().Add(-time.Hour), 1)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(articles)

}
