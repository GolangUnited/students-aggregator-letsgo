package parsers

import "aggregator/models"

type Parser interface {
	Init(map[string]string) error
	ParseAll() ([]models.Article, error)
}
