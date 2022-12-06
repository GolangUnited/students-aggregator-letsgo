package log

import (
	"fmt"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

// Log interface to write messages to log
type Log interface {
	Type() string
	Level() logLevel.LogLevel
	Write(level logLevel.LogLevel, message string, err error) error
	WriteInfo(message string) error
	WriteError(message string, err error) error
}

const (
	timeFormat = "02-01-2006 15:04:05.000000"
)

var (
	printf = fmt.Printf
)

func WriteError(message string, err error) error {
	_, e := printf("%s %s: %v\n", time.Now().Format(timeFormat), message, err)
	return e
}
