package logbot

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	// Настройка логгера при старте приложения
	Log = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true, // Добавляет файл и строку вызова
		}),
	)
}

// WithRequestContext создает логгер с контекстом запроса
func WithRequestContext(requestID string) *slog.Logger {
	return Log.With(
		slog.String("request_id", requestID),
	)
}
