package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

type RotatingLogger struct {
	file     *os.File
	logsDir  string
	FileName string
	exit     chan struct{}
	lock     sync.Mutex // todo
}

func NewRotatingLogger(logDir string) (*RotatingLogger, error) {
	if logDir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		logDir = fmt.Sprintf("%s/logs", dir)
	}
	exists, err := FileExists(logDir)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	fileName := time.Now().Format("2006-01-02.log")
	file, err := os.OpenFile(path.Join(logDir, fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return &RotatingLogger{
		logsDir:  logDir,
		file:     file,
		FileName: fileName,
		exit:     make(chan struct{}, 1),
	}, nil
}

func (rl *RotatingLogger) Flush() error {
	if rl.file != nil {
		return rl.file.Sync()
	}
	return nil
}

func (rl *RotatingLogger) Write(p []byte) (n int, err error) {
	if rl.file != nil {
		return rl.file.Write(p)
	}
	return 0, nil
}
func (rl *RotatingLogger) Exit() error {
	close(rl.exit)
	if rl.file != nil {
		return rl.file.Close()
	}
	return nil
}

func (rl *RotatingLogger) rotate() error {
	daySecond := int64(86400)
	timeLeft := daySecond - time.Now().Unix()%daySecond + 61
	fileTimer := time.After(time.Duration(timeLeft) * time.Second)
	for {
		select {
		case <-fileTimer:
			fileName := time.Now().Format("2006-01-02.log")
			newfile, err := os.OpenFile(path.Join(rl.logsDir, fileName), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				log.Println("open file error:", err)

			}
			tmpFile := rl.file
			rl.file = newfile
			err = tmpFile.Close()
			if err != nil {
				log.Printf("log file close error: %v \n", err)
			}
			rl.FileName = fileName
			timeLeft := daySecond - time.Now().Unix()%daySecond + 61
			fileTimer = time.After(time.Duration(timeLeft) * time.Second)
		case <-rl.exit:
			return nil
		}
	}
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("stat error: %v", err)
}
