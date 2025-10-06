package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger and zap.SugaredLogger
type Logger struct {
	log   *zap.Logger
	sugar *zap.SugaredLogger
}

// New creates a new Logger instance based on the mode
func New(mode string) (*Logger, error) {
	var cfg zap.Config

	switch mode {
	case "dev":
		cfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Development:      true,
			Encoding:         "console",
			EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	case "prod":
		cfg = zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
			Development:      false,
			Encoding:         "json",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	default:
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{
		log:   logger,
		sugar: logger.Sugar(),
	}, nil
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.log.Sync()
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.sugar.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.sugar.Fatalf(format, args...)
}

// InitTestLogger creates a new Logger instance for testing
func InitTestLogger() *Logger {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &Logger{
		log:   logger,
		sugar: logger.Sugar(),
	}
}
