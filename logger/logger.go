package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var logger *zap.Logger
var err error

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
func Fatal(msg string, ctx ...interface{}) {
	logger.Fatal(fmt.Sprintf(msg, ctx...))
}

func InitLogger() error {
	logger, err = newLogger()
	if err != nil {
		return fmt.Errorf("new logger error:%v", err)
	}
	return err
}

func Close() error {
	if logger != nil {
		return logger.Sync()
	}
	return nil
}

func newLogger() (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger, nil
}
