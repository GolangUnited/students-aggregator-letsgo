package logLevel

import (
	"testing"
)

const (
	llUndefined = ""
	llErrors    = "errors"
	llInfo      = "info"
)

func Test(t *testing.T) {
	var l LogLevel

	if l != llUndefined {
		t.Errorf("incorrect log level \"%v\", expected \"%s\" (undefined)", l, llUndefined)
	}

	l = Errors

	if l != llErrors {
		t.Errorf("incorrect log level \"%v\", expected \"%s\" (undefined)", l, llErrors)
	}

	l = Info

	if l != llInfo {
		t.Errorf("incorrect log level \"%v\", expected \"%s\" (undefined)", l, llInfo)
	}
}
