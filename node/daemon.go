package node

import (
	"github.com/lightec-xyz/daemon/logger"
	"os"
	"os/signal"
	"syscall"
)

type IAgent interface {
	Run() error
	Init() error
	Close() error
	Name() string
}

type Daemon struct {
	agents []IAgent
}

func NewDaemon(config Config) *Daemon {
	btcAgent := NewBitcoinAgent()
	ethAgent := NewEthereumAgent()
	return &Daemon{
		agents: []IAgent{btcAgent, ethAgent},
	}
}

func (d *Daemon) Init() error {
	for _, node := range d.agents {
		if err := node.Init(); err != nil {
			logger.Error("%v:init node error %v", node.Name(), err)
			return err
		}
	}
	return nil
}

func (d *Daemon) Close() error {
	for _, node := range d.agents {
		if err := node.Close(); err != nil {
			logger.Error("%v:close node error %v", node.Name(), err)
			//need continue,close next node
		}
	}
	return nil
}

func (d *Daemon) Run() error {
	err := d.Init()
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	ch := make(chan os.Signal, 1)

	for _, node := range d.agents {
		go func(cNode IAgent) {
			err := cNode.Run()
			if err != nil {
				logger.Error("%v run error %v", cNode.Name(), err)
				ch <- syscall.SIGQUIT
			}
		}(node)
	}
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		msg := <-ch
		switch msg {
		case syscall.SIGHUP:
			logger.Info("daemon get SIGHUP")

		case syscall.SIGQUIT:
			fallthrough
		case syscall.SIGTERM:
			logger.Info("get shutdown sigterm, please wait ...")
			err := d.Close()
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}
}
