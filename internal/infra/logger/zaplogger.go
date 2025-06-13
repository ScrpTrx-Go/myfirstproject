package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger(isProd bool, logFilePath string) (*ZapLogger, error) {
	// Уровень логирования
	level := zapcore.DebugLevel

	// Конфигурация вывода в файл
	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level)

	// Конфигурация вывода в консоль
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)

	// Объединяем: и в файл, и в консоль
	combinedCore := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return &ZapLogger{sugar: logger.Sugar()}, nil
}

func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Sync() error {
	return l.sugar.Sync()
}
