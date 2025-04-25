package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Level string

type Logger struct {
	*zap.Logger
	Adapter LoggerAdapter
}

type LoggerAdapter struct {
	lg *zap.Logger
}

const (
	DebugLevel  Level = "debug"
	InfoLevel   Level = "info"
	WarnLevel   Level = "warn"
	DPanicLevel Level = "dpanic"
	PanicLevel  Level = "panic"
	FalalLevel  Level = "fatal"
)

func NewLogger(serviceName string, level Level) *Logger {
	levelsToZap := map[Level]zapcore.Level{
		DebugLevel:  zapcore.DebugLevel,
		InfoLevel:   zapcore.InfoLevel,
		WarnLevel:   zapcore.WarnLevel,
		DPanicLevel: zapcore.DPanicLevel,
		PanicLevel:  zapcore.PanicLevel,
		FalalLevel:  zapcore.FatalLevel,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + serviceName + `_logfile.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), levelsToZap[level]),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   `./log/` + serviceName + `_error.log`,
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
		}), zap.ErrorLevel),
	)
	lg := zap.New(core)
	return &Logger{
		lg,
		LoggerAdapter{lg},
	}
}

func getFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))

	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	return zapFields
}

func (l *LoggerAdapter) Error(msg string, err error, fields map[string]interface{}) {
	allFields := append([]zap.Field{zap.Error(err)}, getFields(fields)...)
	l.lg.Error(msg, allFields...)
}

func (l *LoggerAdapter) Info(msg string, fields map[string]interface{}) {
	l.lg.Info(msg, getFields(fields)...)
}

func (l *LoggerAdapter) Debug(msg string, fields map[string]interface{}) {
	l.lg.Debug(msg, getFields(fields)...)
}

func (l *LoggerAdapter) Trace(_ string, _ map[string]interface{}) {
	return
}

func (l *LoggerAdapter) With(fields map[string]interface{}) LoggerAdapter {
	return LoggerAdapter{
		lg: l.lg.With(getFields(fields)...),
	}
}
