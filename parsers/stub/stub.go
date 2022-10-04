package stub

import (
	"aggregator/models"
	"fmt"
)

type StubParser struct {
	Url string
}

func NewStubParser() *StubParser {
	return &StubParser{}
}

func (p *StubParser) Init(config map[string]string) error {
	p.Url = config["Url"]

	if p.Url == "" {
		return fmt.Errorf("incorrect config")
	}

	return nil
}

func (p *StubParser) ParseAll() ([]models.Article, error) {
	return []models.Article{{Text: fmt.Sprintf("Url = %s", p.Url)}, {Text: "fghwjhcdghwghd"}}, nil
}
