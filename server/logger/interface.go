package logger

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})

}
// logger interface implements zap logger (hides zaplogger )

func NewLogger()Logger{
	return &NewZapLogger(

	)
}
