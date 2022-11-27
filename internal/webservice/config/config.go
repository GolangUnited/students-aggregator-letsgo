package config

import (
	"github.com/indikator/aggregator_lets_go/internal/log"
)

type Config struct {
	Port   uint16     `yaml:"port"`
	Handle string     `yaml:"handle"`
	Log    log.Config `yaml:"log"`
}
