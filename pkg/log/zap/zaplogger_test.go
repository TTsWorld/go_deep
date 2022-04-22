package zap

import (
	"go.uber.org/zap"
	"log"
	"testing"
)

//RedirectStdLog 将标准库的logger输出重定向到 InfoLevel 提供的 logger
//由于 zap 已经处理了调用者注解、时间戳等，它会自动禁用标准库的注解和前缀。
//It returns a function to restore the original prefix and flags and reset the standard library's output to os.Stderr.
func TestRedirectStdLog(t *testing.T) {
	logger := zap.NewExample()
	defer logger.Sync()

	undo := zap.RedirectStdLog(logger)
	defer undo()

	log.Print("redirected standard library")

}
