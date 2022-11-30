package log

import (
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	"time"
)

// Log interface to write messages to log
type Log interface {
	Write(level logLevel.LogLevel, message string, err error) error
	WriteInfo(message string) error
	WriteError(message string, err error) error
}

const (
	timeFormat = "02-01-2006 15:04:05.000000"
)

func WriteError(message string, err error) error {
	_, e := fmt.Printf("%s %s: %v\n", time.Now().Format(timeFormat), message, err)
	return e
}
