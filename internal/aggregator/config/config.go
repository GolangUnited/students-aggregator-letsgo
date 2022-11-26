package config

import (
	"github.com/indikator/aggregator_lets_go/internal/config/logLevel"
)

type Config struct {
	LogLevel logLevel.LogLevel `yaml:"logLevel"`
}
