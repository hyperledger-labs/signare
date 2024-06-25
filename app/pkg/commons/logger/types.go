package logger

import (
	"log/slog"
	"os"
)

// Level represents the importance or severity of a log event.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelDebug Level = "DEBUG"

	// To enable a more user-friendly log output format, set the environment variable ["DEV_LOCAL_MODE"="yes"].
	devLocalModeKey     = "DEV_LOCAL_MODE"
	devLocalModeEnabled = "yes"
)

// logLevelTranslator used to translate from a custom log level to its corresponding slog level.
var logLevelTranslator = map[Level]slog.Level{
	LevelDebug: slog.LevelDebug,
	LevelInfo:  slog.LevelInfo,
	LevelWarn:  slog.LevelWarn,
	LevelError: slog.LevelError,
}

func isLocalMode() bool {
	return os.Getenv(devLocalModeKey) == devLocalModeEnabled
}
