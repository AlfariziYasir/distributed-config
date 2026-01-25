package utils

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type ZapGormLogger struct {
	zapLogger     *zap.Logger
	LogLevel      logger.LogLevel
	SlowThreshold time.Duration
}

func NewZapGormLogger(zapLogger *zap.Logger, logLevel logger.LogLevel, slowThreshold time.Duration) *ZapGormLogger {
	return &ZapGormLogger{
		zapLogger:     zapLogger,
		LogLevel:      logLevel,
		SlowThreshold: slowThreshold,
	}
}

func (l *ZapGormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

func (l *ZapGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.zapLogger.Sugar().Infof(msg, data...)
	}
}

func (l *ZapGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.zapLogger.Sugar().Warnf(msg, data...)
	}
}

func (l *ZapGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.zapLogger.Sugar().Errorf(msg, data...)
	}
}

func (l *ZapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.LogLevel >= logger.Error:
		l.zapLogger.Error("gorm error",
			zap.Error(err),
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.String("file", utils.FileWithLineNum()),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.zapLogger.Warn("gorm slow query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.String("file", utils.FileWithLineNum()),
		)
	case l.LogLevel == logger.Info:
		l.zapLogger.Info("gorm query",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Duration("elapsed", elapsed),
			zap.String("file", utils.FileWithLineNum()),
		)
	}
}
