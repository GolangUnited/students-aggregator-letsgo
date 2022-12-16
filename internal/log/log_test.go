package log

import (
	"fmt"
	"strings"
	"testing"
)

const (
	raiseErrorMessage = "raiseError"
)

var (
	lastMessage string
)

func mockPrintf(format string, a ...any) (n int, err error) {
	if len(a) > 1 && a[1] == raiseErrorMessage {
		return 0, fmt.Errorf(raiseErrorMessage)
	}

	lastMessage = fmt.Sprintf(format, a...)

	return len(lastMessage), nil
}

func TestWriteError(t *testing.T) {
	originPrintf := printf

	defer func() {
		printf = originPrintf
	}()

	printf = mockPrintf

	err := WriteError("123", nil)

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if !strings.Contains(lastMessage, "123") {
		t.Errorf("incorrect message \"%s\", does not contains \"%s\"", lastMessage, "123")
	}

	err = WriteError("123", fmt.Errorf("error123"))

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}

	if !strings.Contains(lastMessage, "123: error123") {
		t.Errorf("incorrect message \"%s\", does not contains \"%s\"", lastMessage, "123: error123")
	}

	err = WriteError(raiseErrorMessage, nil)

	if err == nil {
		t.Error("no error, but it was expected")
	}

	if err.Error() != raiseErrorMessage {
		t.Errorf("unexpected error %v, expected %s", err, raiseErrorMessage)
	}
}
