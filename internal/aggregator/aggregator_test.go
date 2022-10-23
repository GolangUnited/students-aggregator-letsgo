package aggregator

import (
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/config"
)

const (
	configData = `# Project Aggregator YAML
aggregator:
 nothing:

database:
 url: mongodb://localhost:27018/

webservice:
 port: 8080

parsers:
- stub:
   url: https://stub.com`
)

func TestWorkWithStubParser(t *testing.T) {

	c := config.NewConfig()

	err := c.SetData([]byte(configData))

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	err = c.Read()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(c.Parsers) != 1 {
		t.Errorf("incorrect config parsers count %d, expected %d", len(c.Parsers), 1)
	}

	parsers, err := GetParsers(c.Parsers)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(parsers) != 1 {
		t.Errorf("incorrect parsers count %d, expected %d", len(parsers), 1)
	}

	parser := parsers[0]

	articles, err := parser.ParseAll()

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if len(articles) != 3 {
		t.Errorf("incorrect parsers count %d, expected %d", len(parsers), 3)
	}
}
