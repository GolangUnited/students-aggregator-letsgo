package webservice

import (
	"github.com/indikator/aggregator_lets_go/internal/config/logLevel"
)

type Config struct {
	Handle   string
	Port     uint16
	LogLevel logLevel.LogLevel
}
