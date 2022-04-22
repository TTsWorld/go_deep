package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

type addr struct {
	IP   string
	Port int
}

type request struct {
	URL    string
	Listen addr
	Remote addr
}

func (a addr) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ip", a.IP)
	enc.AddInt("port", a.Port)
	return nil
}

func (r request) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("url", r.URL)
	zap.Inline(r.Listen).AddTo(enc)
	return enc.AddObject("remote", r.Remote)
}

func TestObject(t *testing.T) {
	logger := zap.NewExample()
	defer logger.Sync()

	req := &request{
		URL:    "/test",
		Listen: addr{"127.0.0.1", 8080},
		Remote: addr{"127.0.0.1", 31200},
	}
	logger.Info("new request, in nested object", zap.Object("req", req))
	logger.Info("new request, inline", zap.Inline(req))

}
