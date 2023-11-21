package logger

import (
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"os"
)

var logger log.Logger
var err error

func Trace(msg string, ctx ...interface{}) {
	logger.Trace(fmt.Sprintf(msg, ctx...))
}

func Debug(msg string, ctx ...interface{}) {
	logger.Debug(fmt.Sprintf(msg, ctx...))
}

func Info(msg string, ctx ...interface{}) {
	logger.Info(fmt.Sprintf(msg, ctx...))
}

func Warn(msg string, ctx ...interface{}) {
	logger.Warn(fmt.Sprintf(msg, ctx...))
}

func Error(msg string, ctx ...interface{}) {
	logger.Error(fmt.Sprintf(msg, ctx...))
}

func InitLogger() error {
	logger, err = newLogger()
	if err != nil {
		return fmt.Errorf("new logger error:%v", err)
	}
	return err
}

func newLogger() (log.Logger, error) {
	//todo
	logger := log.New()
	{
		glog := log.NewGlogHandler(log.StreamHandler(os.Stdout, log.TerminalFormat(false)))
		glog.Verbosity(log.LvlTrace)
		logger.SetHandler(glog)
	}
	return logger, nil
}
