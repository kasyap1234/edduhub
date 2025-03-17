package logger

import "go.uber.org/zap"


type ZapLogger struct {
	zap *zap.Logger

}