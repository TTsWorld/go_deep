package log

import (
	"errors"
	"testing"
)

//zaplog  库测试
var err = errors.New("123")

func TestLog(t *testing.T) {
	appname, level, fileDir, devModule := "zaplog1", "debug", "./log", false

	// init logger
	InitLog(appname, level, fileDir, devModule)
	GetLogger().Debug("http start success， port " + "8080")
	GetLogger().Info("http start success， port " + "8080")
	GetLogger().Warn("http start success， port " + "8080")
	GetLogger().Error("http start success， port " + "8080")
	GetLogger().Fatal("http start success， port " + "8080" + err.Error())

}
