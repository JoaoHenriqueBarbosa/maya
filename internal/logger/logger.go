package logger

import (
	"fmt"
)

type LogLevel int

const (
	LevelSilent LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

var (
	currentLevel = LevelSilent
	categories   = make(map[string]bool)
)

func init() {
	initConfig()
}

func SetLevel(level LogLevel) {
	currentLevel = level
}

func EnableCategory(category string) {
	categories[category] = true
}

func DisableCategory(category string) {
	delete(categories, category)
}

func shouldLog(level LogLevel, category string) bool {
	if currentLevel == LevelSilent {
		return false
	}
	if level > currentLevel {
		return false
	}
	if len(categories) > 0 && category != "" {
		return categories[category]
	}
	return true
}

func Error(category string, format string, args ...interface{}) {
	if shouldLog(LevelError, category) {
		fmt.Printf("[ERROR][%s] %s\n", category, fmt.Sprintf(format, args...))
	}
}

func Warn(category string, format string, args ...interface{}) {
	if shouldLog(LevelWarn, category) {
		fmt.Printf("[WARN][%s] %s\n", category, fmt.Sprintf(format, args...))
	}
}

func Info(category string, format string, args ...interface{}) {
	if shouldLog(LevelInfo, category) {
		fmt.Printf("[INFO][%s] %s\n", category, fmt.Sprintf(format, args...))
	}
}

func Debug(category string, format string, args ...interface{}) {
	if shouldLog(LevelDebug, category) {
		fmt.Printf("[DEBUG][%s] %s\n", category, fmt.Sprintf(format, args...))
	}
}

func Trace(category string, format string, args ...interface{}) {
	if shouldLog(LevelTrace, category) {
		fmt.Printf("[TRACE][%s] %s\n", category, fmt.Sprintf(format, args...))
	}
}