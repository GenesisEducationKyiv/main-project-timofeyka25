package logger

import (
	"go.uber.org/zap"
	"log"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func NewZapLogger(path string) *ZapLogger {
	logger, err := setupDefaultZapLogger(path)
	if err != nil {
		log.Println(err.Error(), "failed to create new zap logger")
		return nil
	}

	return &ZapLogger{
		logger: logger.Sugar(),
	}
}

func (l *ZapLogger) Debug(msg string) {
	l.logger.Debug(msg)
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg)
}

func setupDefaultZapLogger(path string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{path}
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
