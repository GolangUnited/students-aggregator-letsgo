package aggregator

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/config"
)

type Aggregator struct {
	config config.Config
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Init(config *config.Config) error {
	a.config = *config

	return nil
}

func (a *Aggregator) Execute() error {
	err := a.config.Read()

	if err != nil {
		return err
	}

	parsers, err := GetParsers(a.config.Parsers)

	if err != nil {
		return err
	}

	for _, v := range parsers {
		articles, err := v.ParseAll()

		if err != nil {
			return err
		}

		fmt.Println(articles)
	}

	fmt.Println(a.config)

	return nil
}
