package common

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	"github.com/indikator/aggregator_lets_go/internal/log/std"
	"github.com/indikator/aggregator_lets_go/internal/log/stub"
)

const (
	unknownLogErrorTemplate      = "unknown log type %s"
	unknownLogLevelErrorTemplate = "unknown log level %s"
)

type UnknownLogTypeError struct {
	Text string
}

func (e *UnknownLogTypeError) Error() string {
	return e.Text
}

type UnknownLogLevelError struct {
	Text string
}

func (e *UnknownLogLevelError) Error() string {
	return e.Text
}

var (
	stdNewLog = std.NewLog
)

func GetLog(config log.Config) (log.Log, error) {
	if config.Level != logLevel.Info && config.Level != logLevel.Errors {
		return nil, &UnknownLogLevelError{
			Text: fmt.Sprintf(unknownLogLevelErrorTemplate, config.Level),
		}
	}
	var l log.Log

	switch config.Type {
	case "stub":
		l = stub.NewLog(config.Level)
	case "std":
		l = stdNewLog(config.Level)
	default:
		return nil, &UnknownLogTypeError{
			Text: fmt.Sprintf(unknownLogErrorTemplate, config.Type),
		}
	}

	return l, nil
}
