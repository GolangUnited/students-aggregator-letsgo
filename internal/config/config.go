package config

import (
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database   db.Config
	WebService webservice.Config
	Parsers    []parser.Config
}

func NewConfig() *Config {
	return &Config{}
}

type parserYamlConfig struct {
	Url string `yaml:"url"`
}

type yamlConfig struct {
	Database   db.Config                     `yaml:"database"`
	WebService webservice.Config             `yaml:"webservice"`
	Parsers    []map[string]parserYamlConfig `yaml:"parsers"`
}

func (c *Config) Read(fileName string) (err error) {
	file, err := os.ReadFile(fileName)

	if err != nil {

		return
	}

	var yc yamlConfig

	err = yaml.Unmarshal(file, &yc)

	if err != nil {

		return
	}

	c.Database = yc.Database
	c.WebService = yc.WebService
	c.Parsers = make([]parser.Config, len(yc.Parsers))

	for i, v := range yc.Parsers {
		for name, p := range v {
			c.Parsers[i] = parser.Config{Name: name, Url: p.Url}
		}
	}

	return nil
}
