package logger

import (
	"log"
)

// LogLevel is kind of Loglavel
type LogLevel string

const (
	// LogDebug is
	LogDebug LogLevel = "Debug"
	// LogInfo is
	LogInfo LogLevel = "Info"
	// LogSilent is
	LogSilent LogLevel = "Silent"
)

// Log is
type Log struct {
	Level LogLevel
}

// NewLogger is logger constructor
func NewLog