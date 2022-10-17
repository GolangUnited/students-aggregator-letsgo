package config

import (
	"os"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"gopkg.in/yaml.v3"
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

func (c *Config) ReadFile(fileName string) (err error) {
	file, err := os.ReadFile(fileName)

	if err != nil {
		return
	}

	return c.Read(file)
}

func (c *Config) Read(data []byte) (err error) {
	var yc yamlConfig

	err = yaml.Unmarshal(data, &yc)

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
