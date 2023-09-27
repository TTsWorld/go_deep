package log

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// async zaplog
type ZapLogOper interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

// 同步日志，直接写
type synczaplogger struct {
	zaplog *zap.Logger
}

// 异步日志
type asynczaplogger struct {
	zaplog   *zap.Logger
	masyslog *asynclogger
}

type zaplogger struct {
	zaplog   *zap.Logger
	masyslog *asynclogger
}

func setzaplogger(nwriters map[zapcore.Level]zapcore.WriteSyncer) *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		MessageKey:     "msg",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		//EncodeName: zapcore.FullNameEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	lvlenabler := zap.NewAtomicLevel()
	fcore := newfilecore(encoder, nwriters, lvlenabler)
	zaplog := zap.New(fcore)
	return zaplog
}

func newsynczaplogger(nwriters map[zapcore.Level]zapcore.WriteSyncer) *synczaplogger {

	retlogger := &synczaplogger{}
	retlogger.zaplog = setzaplogger(nwriters)
	return retlogger
}

func newasynczaplogger(nwriters map[zapcore.Level]zapcore.WriteSyncer) *asynczaplogger {
	retlogger := &asynczaplogger{}
	retlogger.zaplog = setzaplogger(nwriters)
	//设置异步的操作
	retlogger.masyslog = newAsyncLogger()
	return retlogger
}

// 返回接口
func NewZapLogger(wtMode int, nwriters map[zapcore.Level]zapcore.WriteSyncer) ZapLogOper {
	if wtMode == 0 {
		return newsynczaplogger(nwriters)
	} else {
		return newasynczaplogger(nwriters)
	}
}

func NewAllZapLogger(nwriters map[zapcore.Level]zapcore.WriteSyncer) ZapLogOper {
	retzaplogger := &zaplogger{}
	retzaplogger.zaplog = setzaplogger(nwriters)
	retzaplogger.masyslog = newAsyncLogger()
	return retzaplogger
}

func (log *zaplogger) Debug(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Debug, msg, fields...)
}

func (log *zaplogger) Info(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Info, msg, fields...)
}

func (log *zaplogger) Warn(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Warn, msg, fields...)
}

func (log *zaplogger) Error(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Error, msg, fields...)
}

func (log *zaplogger) DPanic(msg string, fields ...zap.Field) {
	log.zaplog.DPanic(msg, fields...)
}

func (log *zaplogger) Panic(msg string, fields ...zap.Field) {
	log.zaplog.Panic(msg, fields...)
}

func (log *zaplogger) Fatal(msg string, fields ...zap.Field) {
	log.zaplog.Fatal(msg, fields...)
}

// -------
func (log *synczaplogger) Debug(msg string, fields ...zap.Field) {
	log.zaplog.Debug(msg, fields...)
}

func (log *synczaplogger) Info(msg string, fields ...zap.Field) {
	log.zaplog.Info(msg, fields...)
}

func (log *synczaplogger) Warn(msg string, fields ...zap.Field) {
	log.zaplog.Warn(msg, fields...)
}

func (log *synczaplogger) Error(msg string, fields ...zap.Field) {
	log.zaplog.Error(msg, fields...)
}

func (log *synczaplogger) DPanic(msg string, fields ...zap.Field) {
	log.zaplog.DPanic(msg, fields...)
}

func (log *synczaplogger) Panic(msg string, fields ...zap.Field) {
	log.zaplog.Panic(msg, fields...)
}

func (log *synczaplogger) Fatal(msg string, fields ...zap.Field) {
	log.zaplog.Fatal(msg, fields...)
}

func (log *asynczaplogger) Debug(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Debug, msg, fields...)
}

func (log *asynczaplogger) Info(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Info, msg, fields...)
}

func (log *asynczaplogger) Warn(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Warn, msg, fields...)
}

func (log *asynczaplogger) Error(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Error, msg, fields...)
}

func (log *asynczaplogger) DPanic(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.DPanic, msg, fields...)
}

func (log *asynczaplogger) Panic(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Panic, msg, fields...)
}

func (log *asynczaplogger) Fatal(msg string, fields ...zap.Field) {
	log.masyslog.doAsyncLog(log.zaplog.Fatal, msg, fields...)
}

// 异步的简单实现
const (
	cst_defmsgsize       = 200
	cst_deffieldnums     = 2
	cst_defmaxlogquenums = 500
)

var _asyncMsgPool = sync.Pool{
	New: func() interface{} {
		return &asyncMsg{msg: make([]byte, 0, cst_defmsgsize),
			fields: make([]zap.Field, 0, cst_deffieldnums)}
	},
}

func newAsyncLogger() *asynclogger {
	mret := &asynclogger{logMsgCh: make(chan *asyncMsg, cst_defmaxlogquenums)}
	mret.start()
	return mret
}

func getAsyncMsg() *asyncMsg {
	return _asyncMsgPool.Get().(*asyncMsg)
}

func putAsynMsg(e *asyncMsg) {
	e.msg = e.msg[:0]
	e.fields = e.fields[:0]
	e.f = nil
	_asyncMsgPool.Put(e)
}

type asyncMsg struct {
	msg    []byte
	fields []zap.Field
	f      func(msg string, fields ...zap.Field)
}

// 异步日志
type asynclogger struct {
	//异步队列数据
	logMsgCh chan *asyncMsg
}

func (log *asynclogger) start() {
	go log.doWriteLog()
}

func (log *asynclogger) doWriteLog() {
	for logdata := range log.logMsgCh {
		//可以分池来处理
		logdata.f(string(logdata.msg), logdata.fields...)
		putAsynMsg(logdata)
	}
}

func (log *asynclogger) doAsyncLog(f func(msg string, fields ...zap.Field), msg string, fields ...zap.Field) {
	logdata := getAsyncMsg()
	logdata.msg = append(logdata.msg, msg...)
	logdata.fields = append(logdata.fields, fields...)
	logdata.f = f
	log.logMsgCh <- logdata
}

// zap.Core接口的实现
type filecore struct {
	zapcore.LevelEnabler
	enc     zapcore.Encoder
	writers map[zapcore.Level]zapcore.WriteSyncer
}

type FileCore interface {
	zapcore.Core
}

func newfilecore(enc zapcore.Encoder, nwriters map[zapcore.Level]zapcore.WriteSyncer, enab zapcore.LevelEnabler) FileCore {
	retcore := &filecore{
		LevelEnabler: enab,
		enc:          enc,
	}
	retcore.writers = make(map[zapcore.Level]zapcore.WriteSyncer)
	for k, v := range nwriters {
		retcore.writers[k] = v
	}
	return retcore
}

func (c *filecore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (c *filecore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *filecore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	if _, ok := c.writers[ent.Level]; !ok {
		return errors.New("no set level writers")
	}
	_, err = c.writers[ent.Level].Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		c.writers[ent.Level].Sync()
	}
	return nil
}

func (c *filecore) Sync() error {
	return nil
}

func (c *filecore) clone() *filecore {
	retclone := &filecore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
	}
	retclone.writers = make(map[zapcore.Level]zapcore.WriteSyncer)
	for k, v := range c.writers {
		retclone.writers[k] = v
	}
	return retclone
}
