// internal/config/logger.go
package config

import (
	"log/slog"
	"os"
)

type Component string

const (
	ComponentAPI Component = "api"
	ComponentBot Component = "bot"
)

// NewLogger создаёт логгер с указанным компонентом
func NewLogger(component Component) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo, // можно вынести в env потом
		AddSource: true,           // опционально, в dev удобно
	})

	// Вот здесь магия — добавляем постоянный атрибут
	return slog.New(handler).With("component", component)
}