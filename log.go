package logos

import (
	"log/slog"
	"os"
	"strings"
)

const (
	Debug string = "DEBUG"
	Info  string = "INFO"
	Warn  string = "WARN"
	Error string = "ERROR"
)

type LogEntry struct {
	logger *slog.Logger
	app    map[string]interface{}
}

// NewLogger() creates a modernized subset of the RFC 5424 syslog standard.
// Note: log level is set in the environment variable LOG_LEVEL
//
// Example: set LOG_LEVEL=INFO to see INFO level logs
// | Log call | Level | Result  |
// | -------- | ----- | ------  |
// | DEBUG    | -4    | ❌ skip |
// | INFO     | 0     | ✅ show |
// | WARN     | 4     | ✅ show |
// | ERROR    | 8     | ✅ show |
func NewLogger() *LogEntry {
	return &LogEntry{
		logger: ConfigureLogger(),
		app: map[string]interface{}{
			"source": os.Getenv("APP_NAME"),
		},
	}
}

func ConfigureLogger() *slog.Logger {
	logFormatter := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: GetLogLevel(),
	})
	return slog.New(logFormatter)
}

func (log *LogEntry) appLogger(args ...any) (results []any) {
	for k, v := range log.app {
		results = append(results, k, v)
	}

	countArgs := len(args)
	if countArgs > 0 {
		if countArgs%2 != 0 {
			log.logger.Error("args must be in pairs k/v", args...)
		} else {
			results = append(results, args...)
		}
	}

	return
}

func GetLogLevel() slog.Level {
	level := strings.ToUpper(os.Getenv("LOG_LEVEL"))

	switch level {
	case Debug:
		return slog.LevelDebug
	case Info:
		return slog.LevelInfo
	case Warn:
		return slog.LevelWarn
	case Error:
		return slog.LevelError
	default:
		return slog.LevelInfo // default fallback
	}
}
