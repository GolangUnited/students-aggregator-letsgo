package parser

import (
	"github.com/indikator/aggregator_lets_go/internal/config/logLevel"
)

type Config struct {
	Name     string
	URL      string
	IsLocal  bool
	LogLevel logLevel.LogLevel
}
