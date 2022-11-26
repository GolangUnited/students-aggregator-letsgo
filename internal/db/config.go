package db

import (
	"github.com/indikator/aggregator_lets_go/internal/config/logLevel"
)

type Config struct {
	Name     string            `yaml:"name"`
	Url      string            `yaml:"url"`
	LogLevel logLevel.LogLevel `yaml:"logLevel"`
}
