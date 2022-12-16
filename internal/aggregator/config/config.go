package config

import (
	"github.com/indikator/aggregator_lets_go/internal/log"
)

type Config struct {
	Log log.Config `yaml:"log"`
}
