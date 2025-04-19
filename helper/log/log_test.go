package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	QuickStart()
	logger.Info("test a log")
	logger.Error("test a log")
}
