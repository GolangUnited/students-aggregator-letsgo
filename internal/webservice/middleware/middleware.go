package middleware

import (
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/log"
	"net/http"
)

func LoggingHandler(next http.Handler, l log.Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.WriteInfo(fmt.Sprintf("User %v got last news", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}
