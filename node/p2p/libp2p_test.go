package p2p

import (
	"github.com/lightec-xyz/daemon/logger"
	"testing"
	"time"
)

func TestNewLibP2p01(t *testing.T) {
	err := logger.InitLogger(nil)
	if err != nil {
		t.Error(err)
	}
	node, err := NewLibP2p(NewP2pConfig("", 6001, nil))
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			data := "nihao"
			err := node.Broadcast(&Msg{
				Type: Msg_Hello.Enum(),
				Hello: &Hello{
					Address: &data,
				},
			})
			if err != nil {
				t.Error(err)
			}
		}
	}()
	node.Run()
	ch := make(chan struct{}, 1)
	<-ch
}

func TestNewLibP2p02(t *testing.T) {
	err := logger.InitLogger(nil)
	if err != nil {
		t.Error(err)
	}
	node, err := NewLibP2p(NewP2pConfig("", 6002,
		[]string{""}))
	if err != nil {
		t.Error(err)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			data := "hello"
			err := node.Broadcast(&Msg{
				Type: Msg_Hello.Enum(),
				Hello: &Hello{
					Address: &data,
				},
			})
			if err != nil {
				t.Error(err)
			}
		}
	}()
	node.Run()
	ch := make(chan struct{}, 1)
	<-ch
}
func TestNewLibP2p03(t *testing.T) {
	err := logger.InitLogger(nil)
	if err != nil {
		t.Error(err)
	}
	node, err := NewLibP2p(NewP2pConfig("", 6003,
		[]string{""}))
	if err != nil {
		t.Error(err)
	}
	node.Run()
	ch := make(chan struct{}, 1)
	<-ch
}
