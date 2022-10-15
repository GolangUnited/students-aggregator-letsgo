package main

import (
	"github.com/indikator/aggregator_lets_go/internal/webservice/last_news"
)

func main() {
	handle := "/last_news"

	ws := last_news.NewWebservice(handle)
	last_news.RunServer(ws, handle)
}
