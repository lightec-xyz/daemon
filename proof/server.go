package proof

import (
	"fmt"
	"github.com/lightec-xyz/daemon/common"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc"
)

type Server struct {
	customServer *rpc.CustomServer
	rpcServer    *rpc.Server
	mode         common.Mode
}

func NewServer(url string, mode common.Mode, handler interface{}) (*Server, error) {
	if mode == common.Custom {
		customServer, err := rpc.NewCustomServer(url, handler)
		if err != nil {
			logger.Error("new custom ws server error:%v", err)
			return nil, err
		}
		return &Server{
			customServer: customServer,
			mode:         mode,
		}, nil
	} else if mode == common.Cluster {
		wsServer, err := rpc.NewWsServer("zkbtc", url, handler)
		if err != nil {
			logger.Error("new server error:%v", err)
			return nil, err
		}
		return &Server{
			rpcServer: wsServer,
			mode:      mode,
		}, nil
	} else {
		return nil, fmt.Errorf("unknown mode:%v", mode)
	}
}

func (s *Server) init() error {

	return nil
}

func (s *Server) Run() {
	if s.customServer != nil {
		err := s.customServer.Run()
		if err != nil {
			logger.Error("custom server run error:%v", err)
		}
	}
	if s.rpcServer != nil {
		go s.rpcServer.Run()
	}
}

func (s *Server) Close() error {
	if s.customServer != nil {
		err := s.customServer.Close()
		if err != nil {
			logger.Error("custom server close error:%v", err)
		}
	}
	if s.rpcServer != nil {
		err := s.rpcServer.Shutdown()
		if err != nil {
			logger.Error("rpc server close error:%v", err)
		}
	}
	return nil
}
