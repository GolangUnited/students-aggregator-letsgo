package common

import (
	"fmt"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/std"
	"github.com/indikator/aggregator_lets_go/internal/log/stub"
)

const (
	unknownLogErrorTemplate = "unknown log type %s"
)

type UnknownLogTypeError struct {
	Text string
}

func (e *UnknownLogTypeError) Error() string {
	return e.Text
}

func GetLog(config log.Config) (log.Log, error) {
	var l log.Log

	switch config.Type {
	case "stub":
		l = stub.NewLog(config.Level)
	case "std":
		l = std.NewLog(config.Level)
	default:
		return nil, &UnknownLogTypeError{
			Text: fmt.Sprintf(unknownLogErrorTemplate, config.Type),
		}
	}

	return l, nil
}
