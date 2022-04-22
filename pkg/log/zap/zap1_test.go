package zap

import (
	"go.uber.org/zap"
	"testing"
	"time"
)

// doc:
// https://pkg.go.dev/go.uber.org/zap#section-readme

func TestSugarZap(t *testing.T) {
	url := "http://baidu.com"
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
}

func TestZap(t *testing.T) {
	url := "http://baidu.com"
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	logger.Info("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("url", time.Second),
	)
}

func TestExampleZap(t *testing.T) {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")

}
