package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"testing"
)

// zap 库基本功能 test

func TestZap(t *testing.T) {
	logger, _ = zap.NewProduction()
	logger.Debug("Debug")
	logger.Info("Info")
	logger.Warn("Warn")
	logger.Error("Error")
}

func TestZapSugaredLogger(t *testing.T) {
	logger, _ = zap.NewProduction()
	sugarLogger := logger.Sugar()
	sugarLogger.Debug("Debug")
	sugarLogger.Info("Info")
	sugarLogger.Warn("Warn")
	sugarLogger.Error("Error")
	defer sugarLogger.Sync()
}

// 测试 zaplog 自定义 encoder
// 使用 json 作为输出格式, 输出到文件
func TestZapCustomJsonOut(t *testing.T) {
	fd, _ := os.Create("test.log")
	writereSync := zapcore.AddSync(fd)

	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	core := zapcore.NewCore(encoder, writereSync, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	sugarLogger.Debug("Debug")
	sugarLogger.Info("Info")
	sugarLogger.Warn("Warn")
	sugarLogger.Error("Error")

}

// 测试 zaplog 自定义 encoder
// 输出到控制台
func TestZapCustomConsoleOut(t *testing.T) {
	writereSync := zapcore.AddSync(os.Stdout)

	//设置 encodeconfig
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	encoder := zapcore.NewConsoleEncoder(encoderCfg)
	core := zapcore.NewCore(encoder, writereSync, zapcore.WarnLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	sugarLogger.Debug("Debug")
	sugarLogger.Info("Info")
	sugarLogger.Warn("Warn")
	sugarLogger.Error("Error")
}

// 测试 zaplog 自定义 encoder
// 自动日志切割
func TestZapCustomLogSplit(t *testing.T) {
	writereSync := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writereSync, zapcore.DebugLevel)

	logger := zap.New(core)
	sugarLogger := logger.Sugar()

	for i := 0; i < 1000000; i++ {
		sugarLogger.Debug("Debug")
		sugarLogger.Info("Info")
		sugarLogger.Warn("Warn")
		sugarLogger.Error("Error")
	}
}

// 返回使用 Lumberack 切割日志的 zapcore.WriteSyncer
func getLogWriter2() zapcore.WriteSyncer {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberjackLogger)
}

func getEncoder2() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
func TestTest(t *testing.T) {
}
