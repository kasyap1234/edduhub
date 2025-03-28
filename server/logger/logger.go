package logger

import "go.uber.org/zap"

type ZapLogger struct {
	Logger *zap.Logger
}

func NewZapLogger() Logger {
	logger, _ := zap.NewProduction()
	return &ZapLogger{
		Logger: logger,
	}
}

func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.Logger.Info(msg, convertFields(fields...)...)
}
func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.Logger.Error(msg, convertFields(fields)...)

}

func convertFields(fields ...interface{}) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		if f, ok := field.(zap.Field); ok {
			zapFields[i] = f
		}
	}
	return zapFields
}
