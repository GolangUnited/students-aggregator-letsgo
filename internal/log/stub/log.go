package stub

import (
	ilog "github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
)

type log struct {
	level       logLevel.LogLevel
	lastMessage string
	lastError   error
}

// NewLog create an instance of log
func NewLog(level logLevel.LogLevel) ilog.Log {
	return &log{level: level}
}

func (l *log) Type() string {
	return "stub"
}

func (l *log) Level() logLevel.LogLevel {
	return l.level
}

func (l *log) Write(level logLevel.LogLevel, message string, err error) error {
	if l.level == level || l.level == logLevel.Info && level == logLevel.Errors {
		l.lastMessage = message
		l.lastError = err
	}

	return nil
}

func (l *log) WriteInfo(message string) error {
	return l.Write(logLevel.Info, message, nil)
}

func (l *log) WriteError(message string, err error) error {
	return l.Write(logLevel.Errors, message, err)
}
