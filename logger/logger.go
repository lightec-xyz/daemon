package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var logger *Logger

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

func InitLogger(cfg *LogCfg) error {
	lg, err := newLogger(cfg)
	if err != nil {
		return fmt.Errorf("new logger error:%v", err)
	}
	logger = lg
	return nil
}

func Close() error {
	return logger.Close()
}

type LogCfg struct {
	LogDir   string
	IsStdout bool
	File     bool
}

type Logger struct {
	cfg            *LogCfg
	rotatingLogger *lumberjack.Logger
	*zap.Logger
	exit chan struct{}
}

func (l *Logger) Close() error {
	close(l.exit)
	_ = l.Sync()
	if l.rotatingLogger != nil {
		_ = l.rotatingLogger.Close()
	}
	return nil

}

func (l *Logger) rotating() error {
	daySecond := int64(86400)
	timeLeft := daySecond - time.Now().Unix()%daySecond + 61
	fileTimer := time.After(time.Duration(timeLeft) * time.Second)
	for {
		select {
		case <-fileTimer:
			fileName := fmt.Sprintf("%s/%s", l.cfg.LogDir, time.Now().Format("2006-01-02.log"))
			err := l.rotatingLogger.Close()
			if err != nil {
				fmt.Printf("log rotating error: %v %v \n", fileName, err)
				continue
			}
			l.rotatingLogger = newRotatingLogger(fileName)
			timeLeft := daySecond - time.Now().Unix()%daySecond + 61
			fileTimer = time.After(time.Duration(timeLeft) * time.Second)
		case <-l.exit:
			return nil
		}
	}
}

func newLogger(cfg *LogCfg) (*Logger, error) {
	if cfg == nil {
		cfg = defaultLogCfg()
	}
	var writeSyncers []zapcore.WriteSyncer
	var rotatingLogger *lumberjack.Logger
	var encoderConfig = newStdEncCfg()
	if cfg.File {
		fileName := fmt.Sprintf("%s/%s", cfg.LogDir, time.Now().Format("2006-01-02.log"))
		rotatingLogger = newRotatingLogger(fileName)
		//encoderConfig = newFileEncCfg()
		writeSyncers = append(writeSyncers, zapcore.AddSync(rotatingLogger))
	}
	if cfg.IsStdout {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writeSyncers...),
		zap.DebugLevel,
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l := &Logger{
		Logger:         logger,
		rotatingLogger: rotatingLogger,
		cfg:            cfg,
		exit:           make(chan struct{}, 1),
	}
	if cfg.File {
		go l.rotating()
	}
	return l, nil
}
func newRotatingLogger(fileName string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     2,
		Compress:   false,
	}
}

func newStdEncCfg() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeLevel:  zapcore.LowercaseColorLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
}
func newFileEncCfg() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
	}
}

func defaultLogCfg() *LogCfg {
	return &LogCfg{
		LogDir:   "logs",
		IsStdout: true,
		File:     false,
	}
}
