package logger

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// weiwei's log

var logger *zap.Logger

const (
	loggerKey = iota

	DefaultCode = 1
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

func SetContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(loggerKey), WithContext(ctx).With(fields...))
}

func WithContext(ctx *gin.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	l, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := l.(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return logger
}

func NewBufWriter(n int, writer io.Writer) bufwriter {
	bufw := bufwriter{}
	bufw.bs = make(chan []byte, n)
	go func() {
		defer func() {
			if p := recover(); p != nil {
				logger.Error(fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack()))))
			}
		}()

		for p := range bufw.bs {
			_, _ = writer.Write(p)
		}
	}()
	return bufw
}

func InitWithConfig(level string, filename string) {
	if logger != nil {
		return
	}
	fmt.Println("init logger")
	jack := lumberjack.Logger{
		Filename: filename,
		MaxSize:  1000, // megabytes
		MaxAge:   7,    //days
		Compress: true, // disabled by default
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
		LevelKey:       "x_level",
		TimeKey:        "x_date",
		NameKey:        "x_name",
		CallerKey:      "x_source",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}), NewBufWriter(1000*1000, &jack), zapLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getField(a ...zap.Field) []zap.Field {
	fields := []zap.Field{
		zap.String("x_server", os.Getenv("x_server")),
		zap.Int64("t", time.Now().UnixNano()/1000),
	}
	for _, f := range a {
		fields = append(fields, f)
	}
	return fields
}

func Debugln(args ...interface{}) {
	logger.Debug("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func DebuglnCtx(c *gin.Context, args ...interface{}) {
	WithContext(c).Debug("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func Infoln(args ...interface{}) {
	logger.Info("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func InfolnCtx(c *gin.Context, args ...interface{}) {
	WithContext(c).Info("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func Errorln(args ...interface{}) {
	logger.Error("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func ErrorlnCtx(c *gin.Context, args ...interface{}) {
	WithContext(c).Error("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprint(args...)))...,
	)
}

func Debug(format string, a ...interface{}) {
	logger.Debug("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func DebugCtx(c *gin.Context, format string, a ...interface{}) {
	WithContext(c).Debug("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func Info(format string, a ...interface{}) {
	logger.Info("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func InfoCtx(c *gin.Context, format string, a ...interface{}) {
	WithContext(c).Info("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func Error(format string, a ...interface{}) {
	logger.Error("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func ErrorCtx(c *gin.Context, format string, a ...interface{}) {
	WithContext(c).Error("", getField(
		zap.Int("x_code", DefaultCode),
		zap.Any("x_data", fmt.Sprintf(format, a...)))...,
	)
}

func DebugJson(code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		logger.Debug("", fields...)
	} else {
		logger.Debug("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...,
		)
	}
}

func DebugJsonCtx(c *gin.Context, code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		WithContext(c).Debug("", fields...)
	} else {
		WithContext(c).Debug("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...)
	}
}

func InfoJson(code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		logger.Info("", fields...)
	} else {
		logger.Info("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...)
	}
}

func InfoJsonCtx(c *gin.Context, code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		WithContext(c).Info("", fields...)
	} else {
		WithContext(c).Info("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...)
	}
}

func ErrorJson(code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		logger.Error("", fields...)
	} else {
		logger.Error("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...)
	}
}

func ErrorJsonCtx(c *gin.Context, code int, a interface{}) {
	if m, ok := a.(JsonFormat); ok {
		fields := getField(
			zap.Int("x_code", code),
		)
		for k, v := range m {
			fields = append(fields, zap.Any(k, v))
		}
		WithContext(c).Error("", fields...)
	} else {
		WithContext(c).Error("", getField(
			zap.Int("x_code", code),
			zap.Any("x_data", a))...)
	}
}
