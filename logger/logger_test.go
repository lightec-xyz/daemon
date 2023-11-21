package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	InitLogger()
	Trace("trace %s %v", "test", 13)
	Debug("debug %s %v", "nihao", 11)

}
