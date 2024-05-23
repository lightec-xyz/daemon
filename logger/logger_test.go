package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	err := InitLogger(nil)
	if err != nil {
		t.Error(err)
		return
	}
	Info("info %s %v", "nihao", 11)
	Warn("warn %s %v", "nihao", 11)
	Error("error %s %v", "nihao", 11)
	Debug("debug %s %v", "nihao", 11)
}

func TestLoggerWithFile(t *testing.T) {
	err := InitLogger(&LogCfg{
		LogDir:   "logsDir",
		IsStdout: true,
		File:     true,
	})
	if err != nil {
		t.Error(err)
		return
	}
	Info("info %s %v", "nihao", 11)
	Warn("warn %s %v", "nihao", 11)
	Error("error %s %v", "nihao", 11)
	Debug("debug %s %v", "nihao", 11)
}
