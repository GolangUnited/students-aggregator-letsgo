package logLevel

type LogLevel string

const (
	Undefined LogLevel = ""
	Errors    LogLevel = "errors"
	Info      LogLevel = "info"
)
