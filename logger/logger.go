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
	if logger.discordHook != nil {
		logger.discordHook.Send(fmt.Sprintf(msg, ctx...))
	}
}
func Fatal(msg string, ctx ...interface{}) {
	logger.Fatal(fmt.Sprintf(msg, ctx...))
}

func InitLogger(cfg *LogCfg) error {
	if logger != nil {
		return nil
	}
	lg, err := newLogger(cfg)
	if err != nil {
		return fmt.Errorf("new logger error:%v", err)
	}
	logger = lg
	return nil
}

func Close() error {
	if logger != nil {
		return logger.Close()
	}
	return nil
}

type LogCfg struct {
	LogDir         string
	IsStdout       bool
	File           bool
	DiscordHookUrl string
}

type Logger struct {
	cfg            *LogCfg
	rotatingLogger *RotatingLogger
	*zap.Logger
	discordHook *Discord
}

func (l *Logger) Close() error {
	_ = l.Sync()
	if l.discordHook != nil {
		_, _ = l.discordHook.Call(newDisMsg("discord web hook exit now ...."))
		_ = l.discordHook.Close()
	}
	if l.rotatingLogger != nil {
		_ = l.rotatingLogger.Exit()
	}
	return nil

}

func newLogger(cfg *LogCfg) (*Logger, error) {
	if cfg == nil {
		cfg = defaultLogCfg()
	}
	var writeSyncers []zapcore.WriteSyncer
	var rotatingLogger *RotatingLogger
	var encoderConfig = newStdEncCfg()
	var err error
	if cfg.File {
		rotatingLogger, err = NewRotatingLogger(cfg.LogDir)
		if err != nil {
			return nil, err
		}
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
	if rotatingLogger != nil {
		go rotatingLogger.rotate() // todo
	}
	lg := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	if cfg.DiscordHookUrl != "" {
		discord := NewDiscord(cfg.DiscordHookUrl)
		_, err := discord.Call(newDisMsg("discord web hook running ...."))
		if err != nil {
			return nil, err
		}
		return &Logger{
			Logger:         lg,
			rotatingLogger: rotatingLogger,
			cfg:            cfg,
			discordHook:    discord,
		}, nil
	} else {
		return &Logger{
			Logger:         lg,
			rotatingLogger: rotatingLogger,
			cfg:            cfg,
		}, err
	}

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
		LogDir:   "./logs",
		IsStdout: true,
		File:     false,
	}
}
