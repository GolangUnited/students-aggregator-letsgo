package common

import (
	"errors"
	"fmt"
	"testing"

	"github.com/indikator/aggregator_lets_go/internal/log"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	slog "github.com/indikator/aggregator_lets_go/internal/log/stub"
)

func TestGetLogCorrectLog(t *testing.T) {
	c := log.Config{
		Type:  "stub",
		Level: logLevel.Errors,
	}
	l, err := GetLog(c)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if l.Type() != c.Type {
		t.Errorf("log type %s incorrect, expected %s", l.Type(), c.Type)
	}

	if l.Level() != c.Level {
		t.Errorf("log level %v incorrect, expected %v", l.Level(), c.Level)
	}
}

func TestGetLogIncorrectLog(t *testing.T) {
	c := log.Config{
		Type:  "mock",
		Level: logLevel.Errors,
	}
	_, err := GetLog(c)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unknownLogTypeError *UnknownLogTypeError

	switch {
	case errors.As(err, &unknownLogTypeError):
	default:
		t.Errorf("unexpected error %v", err)
	}

	got := err.Error()
	expected := fmt.Sprintf(unknownLogErrorTemplate, c.Type)

	if got != expected {
		t.Errorf("error message \"%s\" incorrect expected \"%s\"", got, expected)
	}
}

func TestGetLogIncorrectLogLevel(t *testing.T) {
	c := log.Config{
		Type:  "mock",
		Level: logLevel.Undefined,
	}
	_, err := GetLog(c)

	if err == nil {
		t.Error("expect error is missing")
	}

	var unknownLogLevelError *UnknownLogLevelError

	switch {
	case errors.As(err, &unknownLogLevelError):
	default:
		t.Errorf("unexpected error %v", err)
	}

	got := err.Error()
	expected := fmt.Sprintf(unknownLogLevelErrorTemplate, logLevel.Undefined)

	if got != expected {
		t.Errorf("error message \"%s\" incorrect expected \"%s\"", got, expected)
	}
}

func TestGetLogStd(t *testing.T) {
	originNewLog := stdNewLog

	defer func() {
		stdNewLog = originNewLog
	}()

	stdNewLog = slog.NewLog

	c := log.Config{
		Type:  "std",
		Level: logLevel.Info,
	}
	_, err := GetLog(c)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}
