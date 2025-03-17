package logger

import "go.uber.org/zap"


type ZapLogger struct {
	Logger *zap.Logger

}
// logger, _ := zap.NewProduction()
// defer logger.Sync()
// logger.Info("failed to fetch URL",
//   // Structured context as strongly typed Field values.
//   zap.String("url", url),
//   zap.Int("attempt", 3),
//   zap.Duration("backoff", time.Second),
// )

func NewZapLogger()*ZapLogger{
	logger,_:=zap.NewProduction()
 return &ZapLogger{
	Logger: logger,
 }
}

func(l*ZapLogger)Info(msg string,fields...interface{}){
	l.Logger.Info(msg,convertFields(fields...)...) 
}
func(l*ZapLogger)Error(msg string,fields...interface{}){
	l.Logger.Error(msg,convertFields(fields)...)
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

