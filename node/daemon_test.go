package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func TestDaemon(t *testing.T) {
	config := devDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Error(err)
	}
	defer daemon.Close()
	err = daemon.Run()
	if err != nil {
		t.Error(err)
	}
}
func TestDaemon_Demo(t *testing.T) {
	logger.InitLogger()
	logger.Debug("sdfasdfsd")
}
