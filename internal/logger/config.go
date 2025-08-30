//go:build !wasm
// +build !wasm

package logger

import (
	"os"
	"strings"
)

func initConfig() {
	levelStr := os.Getenv("MAYA_LOG_LEVEL")
	switch strings.ToLower(levelStr) {
	case "error":
		currentLevel = LevelError
	case "warn":
		currentLevel = LevelWarn
	case "info":
		currentLevel = LevelInfo
	case "debug":
		currentLevel = LevelDebug
	case "trace":
		currentLevel = LevelTrace
	default:
		currentLevel = LevelSilent
	}

	if catStr := os.Getenv("MAYA_LOG_CATEGORIES"); catStr != "" {
		for _, cat := range strings.Split(catStr, ",") {
			cat = strings.TrimSpace(strings.ToUpper(cat))
			if cat != "" {
				categories[cat] = true
			}
		}
	}
}