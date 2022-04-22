package log

import (
	"go.uber.org/zap/zapcore"
)

const (
	//这里对 zapcore 包的日志级别进行封装，用户可以直接通过当前 log 包设置日志级别，从而
	//隐藏 zap 的实现
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel logs are particularly important errors. In development the
	// logger panics after writing the message.
	//  DPanicLevel = Development PanicLevel
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel logs a message, then panics.
	PanicLevel = zapcore.PanicLevel
	// FatalLevel logs a message, then calls os.Exit(1).
	FatalLevel = zapcore.FatalLevel
)
