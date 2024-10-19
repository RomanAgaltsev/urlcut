package logger

import (
	"log/slog"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

func Initialize() error {
	// Создаем параметры енкодинга
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	// Создаме параметры логера
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	// Пробуем собрать логер
	logger, err := config.Build()
	// Проверяем наличие ошибки
	if err != nil {
		// Есть ошибка, возвращаем nil
		return err
	}
	// Откладываем sync логера
	defer func() { _ = logger.Sync() }()
	// Устанавливаем новый логер slog с бэкендом zap в качестве логера по умолчанию
	slog.SetDefault(slog.New(zapslog.NewHandler(logger.Core(), nil)))
	return nil
}
