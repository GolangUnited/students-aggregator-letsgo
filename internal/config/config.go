package config

import (
	"fmt"
	"os"

	aggregator "github.com/indikator/aggregator_lets_go/internal/aggregator/config"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/internal/parser"
	"github.com/indikator/aggregator_lets_go/internal/webservice"
	"gopkg.in/yaml.v3"
)

type Config struct {
	data       []byte
	Aggregator aggregator.Config
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

func (c *Config) SetDataFromFile(fileName string) error {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	return c.SetData(data)
}

func (c *Config) SetData(data []byte) error {
	c.data = data

	return nil
}

func (c *Config) Read() (err error) {
	if len(c.data) == 0 {
		return fmt.Errorf("data to read not found")
	}

	var yc yamlConfig

	err = yaml.Unmarshal(c.data, &yc)

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
