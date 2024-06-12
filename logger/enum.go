package logger

import "strings"

type LogLevel string

const (
	LogLevel_Debug LogLevel = "debug" // -4
	LogLevel_Info  LogLevel = "info"  // 0
	LogLevel_Warn  LogLevel = "warn"  // 4
	LogLevel_Error LogLevel = "error" // 8
)

var logLevelMap = map[string]LogLevel{
	"debug": LogLevel_Debug,
	"info":  LogLevel_Info,
	"warn":  LogLevel_Warn,
	"error": LogLevel_Error,
}

func (l LogLevel) String() string {
	return string(l)
}

func StringToLevel(level string) LogLevel {
	lower := strings.ToLower(level)
	logLevel := logLevelMap[lower]
	return logLevel
}
