package log

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

// Log interface to write messages to log
type Log interface {
	Write(level logLevel.LogLevel, message string, err error) error
	WriteInfo(message string) error
	WriteError(message string, err error) error
}

func WriteError(message string, err error) error {
	_, e := fmt.Printf("%s: %v\n", message, err)
	return e
}
