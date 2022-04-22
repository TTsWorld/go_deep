package zap

import (
	"go.uber.org/zap"
	"testing"
)

var logger *zap.Logger

func TestCaller(t *testing.T) {
	logger := NewLogger()
	logger.Info("Test msg Main")
	f1(logger)
}

func f1(logger *zap.Logger) {
	logger.Info("Test msg TestFunc")
}

func NewLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.FunctionKey = "func"
	opts := []zap.Option{
		zap.AddCallerSkip(2), // traverse call depth for more useful log lines
		zap.AddCaller(),
	}

	logger, _ = config.Build(opts...)
	return logger
}
