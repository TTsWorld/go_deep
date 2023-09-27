package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//伟伟的 log 封装

var (
	logger2 *zap.Logger
	json    = jsoniter.ConfigCompatibleWithStandardLibrary
)

type key int

const (
	loggerKey key = iota
	maxSize       = 1000
	maxAge        = 7
	bufSize       = 1000 * 1000
)

type JsonFormat map[string]interface{}

type bufwriter struct {
	bs chan []byte
}

func (bw bufwriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p) // must copy, avoid use same slice
	bw.bs <- buf
	return len(p), nil
}

func (bw bufwriter) Sync() error {
	return nil
}

func SetContext(ctx context.Context, fields ...zapcore.Field) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields...))
}

func WithContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger2
	}
	ctxLogger, ok := ctx.Value(loggerKey).(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return logger2
}

func NewBufWriter(n int, writer io.Writer) bufwriter {
	bufw := bufwriter{}
	bufw.bs = make(chan []byte, n)
	go func() {
		for p := range bufw.bs {
			_, _ = writer.Write(p)
		}
	}()
	return bufw
}

func InitWithConfig(level string, filename string) {
	if logger2 != nil {
		return
	}
	jack := lumberjack.Logger{
		Filename: filename,
		MaxSize:  maxSize, // megabytes
		MaxAge:   maxAge,  //days
		Compress: true,    // disabled by default
	}

	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
		break
	case "info":
		zapLevel = zapcore.InfoLevel
		break
	default:
		zapLevel = zapcore.WarnLevel
	}

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "date",
		NameKey:        "name",
		CallerKey:      "source",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}), NewBufWriter(bufSize, &jack), zapLevel)
	logger2 = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getField(a ...zap.Field) []zap.Field {
	fields := []zap.Field{
		zap.String("server", os.Getenv("server")),
		zap.Int64("timestamp", time.Now().Unix()),
		zap.String("event", "mediaReq"),
	}
	for _, f := range a {
		fields = append(fields, f)
	}
	return fields
}

func Debugln(args ...interface{}) {
	logger2.Debug("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func Debugf(format string, args ...interface{}) {
	logger2.Debug("", getField(zap.Any("data", fmt.Sprintf(format, args...)))...)
}

func DebuglnCtx(c context.Context, args ...interface{}) {
	WithContext(c).Debug("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func DebugJson(args interface{}) {
	data, _ := json.Marshal(args)
	logger2.Debug("", getField(zap.Any("data", bytes2string(data)))...)
}

func DebugJsonCtx(c context.Context, args interface{}) {
	data, _ := json.Marshal(args)
	WithContext(c).Debug("", getField(zap.Any("data", bytes2string(data)))...)
}

func Infof(format string, args ...interface{}) {
	logger2.Info("", getField(zap.Any("data", fmt.Sprintf(format, args...)))...)
}

func Infoln(args ...interface{}) {
	logger2.Info("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func InfolnCtx(c context.Context, args ...interface{}) {
	WithContext(c).Info("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func InfoJson(args interface{}) {
	data, _ := json.Marshal(args)
	logger2.Info("", getField(zap.Any("data", bytes2string(data)))...)
}

func InfoField(a ...zap.Field) {
	logger2.Info("", getField(a...)...)
}

func InfoJsonCtx(c context.Context, args interface{}) {
	data, _ := json.Marshal(args)
	WithContext(c).Info("", getField(zap.Any("data", bytes2string(data)))...)
}

func Warnf(format string, args ...interface{}) {
	logger2.Warn("", getField(zap.Any("data", fmt.Sprintf(format, args...)))...)
}

func Warnln(args ...interface{}) {
	logger2.Warn("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func WarnlnCtx(c context.Context, args ...interface{}) {
	WithContext(c).Warn("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func WarnJson(args interface{}) {
	data, _ := json.Marshal(args)
	logger2.Warn("", getField(zap.Any("data", bytes2string(data)))...)
}

func WarnJsonCtx(c context.Context, args interface{}) {
	data, _ := json.Marshal(args)
	WithContext(c).Warn("", getField(zap.Any("data", bytes2string(data)))...)
}

func Errorf(format string, args ...interface{}) {
	logger2.Error("", getField(zap.Any("data", fmt.Sprintf(format, args...)))...)
}

func Errorln(args ...interface{}) {
	logger2.Error("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func ErrorlnCtx(c context.Context, args ...interface{}) {
	WithContext(c).Error("", getField(zap.Any("data", fmt.Sprint(args...)))...)
}

func ErrorJson(args interface{}) {
	data, _ := json.Marshal(args)
	logger2.Error("", getField(zap.Any("data", bytes2string(data)))...)
}

func ErrorJsonCtx(c context.Context, args interface{}) {
	data, _ := json.Marshal(args)
	WithContext(c).Error("", getField(zap.Any("data", bytes2string(data)))...)
}

func bytes2string(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&sh))
}
