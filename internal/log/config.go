package log

import (
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

type Config struct {
	Type  string            `yaml:"type"`
	Level logLevel.LogLevel `yaml:"level"`
}
