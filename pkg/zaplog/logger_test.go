package log

import (
	"errors"
	"testing"
)

var err1 = errors.New("123")
var DefaultCode = 2

func TestLogger(t *testing.T) {
	InitWithConfig("debug", "test.log")
	//Debugln("hello:", "logger")
	//Debug("hello:%s", "logger")
	//DebugJson(0, struct {
	//	AppID string `json:"appid"`
	//}{AppID: "123"})
	//ErrorJson(DefaultCode, JsonFormat{
	//	"error": err.Error(),
	//})
	//ErrorJsonCtx(nil, DefaultCode, JsonFormat{
	//	"errorctx": err.Error(),
	//})
	//InfoJson(DefaultCode, JsonFormat{
	//	"infoJson": err.Error(),
	//})
	//DebugJson(DefaultCode, JsonFormat{
	//	"debugJson": err.Error(),
	//})
}
