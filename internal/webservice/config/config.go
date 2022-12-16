package config

import (
	"github.com/indikator/aggregator_lets_go/internal/log"
)

type Config struct {
	Handle  string     `yaml:"handle"`
	Port    uint16     `yaml:"port"`
	Swagger bool       `yaml:"swagger"`
	Log     log.Config `yaml:"log"`
}
