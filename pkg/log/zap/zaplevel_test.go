package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

func TestAtomicLevel(t *testing.T) {

	//AtomicLevel 是一个具备原子性、可变、动态的logging级别。它允许您在运行时安全地更改logger树
	//（root logger和通过添加上下文创建的任何子logger）的日志级别。
	atom := zap.NewAtomicLevel()

	// To keep the example deterministic, disable timestamps in the output.
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()

	logger.Info("info logging enabled")

	//可在运行时动态的改变日志 level
	atom.SetLevel(zap.ErrorLevel)
	logger.Info("info logging disabled")

}
