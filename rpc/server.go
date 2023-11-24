package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/logger"
	"log"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(addr string, handler interface{}) (*Server, error) {
	//todo
	rpcServ := rpc.NewServer()
	err := rpcServ.RegisterName("zkbtc", handler)
	if err != nil {
		log.Fatal(err)
	}
	rpcServ.SetBatchLimits(100, 100)
	server := &http.Server{
		Addr:    addr,
		Handler: rpcServ,
	}
	http.Handle("/", rpcServ)
	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()
	return &Server{server: server}, nil
}

func (s *Server) Shutdown() error {
	if s.server != nil {
		err := s.server.Shutdown(context.TODO())
		if err != nil {
			logger.Error("server shutdown error:%v", err)
		}
	}
	return nil
}
