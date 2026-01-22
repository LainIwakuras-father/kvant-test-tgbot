package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config — структура со всеми нужными переменными (можно расширять)
type Config struct {
	BotToken   string
	SecretKey  string
	Port       string // например, для сервера
	DebugMode  bool
}

// Load загружает конфигурацию и возвращает ошибку, если что-то критично не задано
func Load() (*Config, error) {
	// Пытаемся загрузить .env (только в dev/local, в prod обычно не нужен)
	if err := godotenv.Load(); err != nil {
		slog.Warn("Не удалось загрузить .env файл (возможно, это production)", "err", err)
	}

	cfg := &Config{}

	// BOT_TOKEN — обязательный
	cfg.BotToken = os.Getenv("BOT_TOKEN")
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN обязателен")
	}

	// SECRET_KEY — обязательный
	cfg.SecretKey = os.Getenv("SECRET_KEY")
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("SECRET_KEY обязателен")
	}

	// PORT — с дефолтом
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = "8080"
		slog.Info("PORT не задан → используется дефолт", "port", cfg.Port)
	}
	// DEBUG_MODE — опционально, по умолчанию false
	if os.Getenv("DEBUG") == "true" {
		cfg.DebugMode = true
	} else {
		cfg.DebugMode = false
	}

	// Можно добавить другие переменные + валидацию здесь

	slog.Info("Конфигурация успешно загружена",
		"bot_token_present", cfg.BotToken != "",
		"secret_key_present", cfg.SecretKey != "",
		"port", cfg.Port,
		"debug_mode", cfg.DebugMode,
	)

	return cfg, nil
}