package main

import (
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
)

const (
	configFilePath = "./configs/config.yaml"
)

func main() {

	handle := "/last_news"

	ws := last_news.NewWebservice(handle)

	cfg := config.NewConfig()
	cfg.SetDataFromFile(configFilePath)

	last_news.RunServer(ws, *cfg, handle)
}
