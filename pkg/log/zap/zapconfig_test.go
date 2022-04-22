package zap

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"testing"
)

func TestAdvanceConfig(t *testing.T) {
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	kafkaEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(kafkaEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(kafkaEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	defer logger.Sync()
	logger.Info("constructed a logger")
}

func TestBasicConfig(t *testing.T) {
	// For some users, the presets offered by the NewProduction, NewDevelopment,
	// and NewExample constructors won't be appropriate. For most of those
	// users, the bundled Config struct offers the right balance of flexibility
	// and convenience. (For more complex needs, see the AdvancedConfiguration
	// example.)
	//
	// See the documentation for Config and zapcore.EncoderConfig for all the
	// available options.
	rawJSON := []byte(`{
	  "level": "debug",
	  "encoding": "json",
	  "outputPaths": ["stdout", "/tmp/logs"],
	  "errorOutputPaths": ["stderr"],
	  "initialFields": {"foo": "bar"},
	  "encoderConfig": {
	    "messageKey": "message",
	    "levelKey": "level",
	    "levelEncoder": "lowercase"
	  }
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		panic(err)
	}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("logger construction succeeded")
}

func TestCheck(t *testing.T) {
	logger := zap.NewExample()
	defer logger.Sync()

	if ce := logger.Check(zap.DebugLevel, "debugging"); ce != nil {
		// If debug-level log output isn't enabled or if zap's sampling would have
		// dropped this log entry, we don't allocate the slice that holds these
		// fields.
		ce.Write(
			zap.String("foo", "bar"),
			zap.String("baz", "quux"),
		)
	}

}
