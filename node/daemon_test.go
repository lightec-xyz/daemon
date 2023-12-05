package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
)

func TestDaemon(t *testing.T) {
	config := localDevDaemonConfig()
	daemon, err := NewDaemon(config)
	if err != nil {
		t.Fatal(err)
	}
	defer daemon.Close()
	err = daemon.Init()
	if err != nil {
		t.Fatal(err)
	}
	err = daemon.Run()
	if err != nil {
		t.Fatal(err)
	}
}
func TestDaemon_Demo(t *testing.T) {
	logger.InitLogger()
	logger.Debug("sdfasdfsd")
}
