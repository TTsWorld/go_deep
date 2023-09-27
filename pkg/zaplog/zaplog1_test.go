package log

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

// zaplogger  库测试
var err = errors.New("123")
var appname, level, fileDir, devModule = "zaplogger", "debug", "./logger", true

// init loggerger
func init() {
	InitLog(appname, level, fileDir, devModule)
}

func TestZapLog(t *testing.T) {
	logger.Debug("http start success， port " + "8080")
	logger.Info("http start success， port " + "8080")
	logger.Warn("http start success， port "+"8080", zap.String("err", err.Error()))
	logger.Error("http start success， port " + "8080")
	//logger.Fatal("http start success， port " + "8080" + err.Error())

}

func TestLoggerRotate(t *testing.T) {
	logger.Debug("http start success， port " + "8080")
	logger.Info("http start success， port " + "8080")
	logger.Warn("http start success， port " + "8080")
	logger.Error("http start success， port " + "8080")
	logger.Fatal("http start success， port " + "8080" + err.Error())

}
