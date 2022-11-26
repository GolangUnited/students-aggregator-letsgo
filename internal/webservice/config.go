package webservice

import (
	"github.com/indikator/aggregator_lets_go/internal/config/logLevel"
)

type Config struct {
	Port     uint16            `yaml:"port"`
	Handle   string            `yaml:"handle"`
	LogLevel logLevel.LogLevel `yaml:"logLevel"`
}
