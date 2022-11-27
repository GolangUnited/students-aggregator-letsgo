package std

import (
	"fmt"

	ilog "github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

type log struct {
	level logLevel.LogLevel
}

// NewLog create an instance of log
func NewLog(level logLevel.LogLevel) ilog.Log {
	return &log{level: level}
}

func (l *log) Write(level logLevel.LogLevel, message string, err error) error {
	if l.level == level || l.level == logLevel.Info && level == logLevel.Errors {
		if err != nil {
			message = fmt.Sprintf("%s: %v", message, err)
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
