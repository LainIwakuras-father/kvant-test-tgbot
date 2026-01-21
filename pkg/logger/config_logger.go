package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const (
	loggerKey ctxKey = "logger"
	// Коды цветов ANSI
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorReset  = "\033[0m"
)

// LogFormatter кастомный форматтер для линейного вывода
type LogFormatter struct {
	ShowCaller bool
}

// Format форматирует запись лога
func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Цвета для уровней
	levelColor := f.getLevelColor(entry.Level)
	
	// Форматируем caller (файл:строка)
	callerStr := ""
	if f.ShowCaller && entry.HasCaller() {
		// Берем только имя файла
		_, filename := filepath.Split(entry.Caller.File)
		callerStr = colorWhite + filename + ":" + fmt.Sprintf("%d", entry.Caller.Line) + colorReset
	}
	
	// Цвет сообщения
	messageColor := colorCyan
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel {
		messageColor = colorRed
	}
	
	// Форматируем уровень
	levelStr := strings.ToUpper(entry.Level.String())
	if len(levelStr) > 4 {
		levelStr = levelStr[0:4]
	}
	
	// Получаем модуль из полей
	module := ""
	if m, ok := entry.Data["module"]; ok {
		module = fmt.Sprintf("[%s]", m)
		delete(entry.Data, "module") // Удаляем модуль из полей чтобы не дублировался
	}
	
	// Собираем другие поля
	var fields []string
	var errorStr string
	
	for k, v := range entry.Data {
		if k == "error" || k == "err" {
			if err, ok := v.(error); ok {
				errorStr = fmt.Sprintf("ERROR: %v", err)
			} else {
				errorStr = fmt.Sprintf("ERROR: %v", v)
			}
		} else {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
	}
	
	// Собираем итоговую строку
	var parts []string
	
	// 1. Caller
	if callerStr != "" {
		parts = append(parts, callerStr)
	}
	
	// 2. Модуль
	if module != "" {
		parts = append(parts, colorPurple+module+colorReset)
	}
	
	// 3. Уровень
	parts = append(parts, levelColor+levelStr+colorReset)
	
	// 4. Сообщение
	parts = append(parts, messageColor+entry.Message+colorReset)
	
	// 5. Поля
	if len(fields) > 0 {
		parts = append(parts, colorWhite+strings.Join(fields, " ")+colorReset)
	}
	
	// 6. Ошибка
	if errorStr != "" {
		parts = append(parts, colorRed+errorStr+colorReset)
	}
	
	// Добавляем время если включен дебаг
	result := strings.Join(parts, " | ")
	if os.Getenv("DEBUG") == "true" {
		timeStr := entry.Time.Format("15:04:05")
		result = colorWhite + "[" + timeStr + "]" + colorReset + " " + result
	}
	
	return []byte(result + "\n"), nil
}

func (f *LogFormatter) getLevelColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return colorPurple
	case logrus.InfoLevel:
		return colorBlue
	case logrus.WarnLevel:
		return colorYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return colorRed
	default:
		return colorWhite
	}
}

// Init инициализирует логгер
func Init() {
	// Настраиваем уровень логирования
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	
	// Включаем вывод caller для WARN и ERROR
	logrus.SetReportCaller(true)
	
	// Настраиваем форматтер
	formatter := &LogFormatter{
		ShowCaller: os.Getenv("LOG_CALLER") == "true",
	}
	logrus.SetFormatter(formatter)
	
	// Настраиваем вывод
	logrus.SetOutput(os.Stdout)
	
	// Добавляем хук для скрытия caller для INFO и DEBUG
	logrus.AddHook(&callerHook{})
}

// callerHook скрывает caller для INFO и DEBUG если не включено явно
type callerHook struct{}

func (h *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *callerHook) Fire(entry *logrus.Entry) error {
	// Показываем caller только для WARN и ERROR или если явно включено
	if entry.Level <= logrus.InfoLevel && os.Getenv("LOG_CALLER") != "true" {
		entry.Caller = nil
	}
	return nil
}

// NewModuleLogger создает логгер для модуля
func NewModuleLogger(module string) *logrus.Entry {
	return logrus.WithField("module", module)
}

// FromContext возвращает логгер из контекста
func FromContext(ctx context.Context) *logrus.Entry {
	if logger, ok := ctx.Value(loggerKey).(*logrus.Entry); ok && logger != nil {
		return logger
	}
	// Возвращаем стандартный логгер
	return logrus.WithField("module", "default")
}

// NewContextWithLogger создает контекст с логгером
func NewContextWithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// WithFields добавляет поля в логгер из контекста
func WithFields(ctx context.Context, fields logrus.Fields) context.Context {
	logger := FromContext(ctx)
	return NewContextWithLogger(ctx, logger.WithFields(fields))
}

// WithField добавляет поле в логгер из контекста
func WithField(ctx context.Context, key string, value interface{}) context.Context {
	return WithFields(ctx, logrus.Fields{key: value})
}

// Helper функции для удобного логирования

func Debug(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Debug(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Debugf(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Info(args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Infof(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Warn(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Warnf(format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Error(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Errorf(format, args...)
}

// WithCallerInfo добавляет информацию о caller
func WithCallerInfo(ctx context.Context) context.Context {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return ctx
	}
	
	funcName := runtime.FuncForPC(pc).Name()
	_, filename := filepath.Split(file)
	
	return WithFields(ctx, logrus.Fields{
		"caller_file": filename,
		"caller_line": line,
		"caller_func": funcName,
	})
}

// GetCaller возвращает информацию о caller для текущего вызова
func GetCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown", 0
	}
	_, filename := filepath.Split(file)
	return filename, line
}