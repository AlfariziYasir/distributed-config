package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

// New -.
func New(level zapcore.Level, service string, appVersion string, isProd bool) (Logger, error) {
	var (
		l   *zap.Logger
		err error
	)

	if isProd {
		l, err = createProdLog(level, service, appVersion)
		if err != nil {
			return Logger{}, err
		}
	} else {

		l, err = createDevLog(level)
		if err != nil {
			return Logger{}, err
		}
	}

	return Logger{
		logger: l,
	}, nil

}

func createProdLog(level zapcore.Level, service string, appVersion string) (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.InitialFields = map[string]any{
		"service": service,
		"version": appVersion,
	}

	return config.Build(zap.AddCaller())
}

func createDevLog(level zapcore.Level) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = zap.NewAtomicLevelAt(level)
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return config.Build(zap.AddCaller())
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *Logger) Sync() error {
	if l.logger != nil {
		return l.logger.Sync()
	}
	return nil
}

func (l *Logger) GetZapLogger() *zap.Logger {
	return l.logger
}
