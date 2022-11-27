package std

import (
	"fmt"
	"time"

	ilog "github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

type log struct {
	level logLevel.LogLevel
}

const (
	timeFormat = "02-01-2006 15:04:05.000000"
)

// NewLog create an instance of log
func NewLog(level logLevel.LogLevel) ilog.Log {
	return &log{level: level}
}

func (l *log) Write(level logLevel.LogLevel, message string, err error) error {
	if l.level == level || l.level == logLevel.Info && level == logLevel.Errors {
		if err != nil {
			message = fmt.Sprintf("%s %s: %v", time.Now().Format(timeFormat), message, err)
		} else {
			message = fmt.Sprintf("%s %s", time.Now().Format(timeFormat), message)
		}

		_, e := fmt.Println(message)
		return e
	}

	return nil
}

func (l *log) WriteInfo(message string) error {
	return l.Write(logLevel.Info, message, nil)
}

func (l *log) WriteError(message string, err error) error {
	return l.Write(logLevel.Errors, message, err)
}
