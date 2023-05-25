package enum

type LoggerLevel string

const (
	LoggerLevel_Debug LoggerLevel = "debug"
	LoggerLevel_Error LoggerLevel = "error"
	LoggerLevel_Info  LoggerLevel = "info"
)
