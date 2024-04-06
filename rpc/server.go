package rpc

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lightec-xyz/daemon/logger"
	"net"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer *http.Server
	name       string
}

func NewServer(name, addr string, handler interface{}) (*Server, error) {
	isOpen := isPortOpen(addr)
	if isOpen {
		return nil, fmt.Errorf("port is open:%v", addr)
	}
	rpcServer := rpc.NewServer()
	err := rpcServer.RegisterName(name, handler)
	if err != nil {
		logger.Error("register name error:%v %v", name, err)
		return nil, err
	}
	rpcServer.SetBatchLimits(BatchRequestLimit, BatchResponseMaxSize)
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        rpcServer,
		ReadTimeout:    HttpReadTimeOut,
		WriteTimeout:   HttpWriteTimeOut,
		MaxHeaderBytes: MaxHeaderBytes,
		IdleTimeout:    30 * time.Minute,
	}
	return &Server{httpServer: httpServer, name: name}, nil
}

func NewWsServer(name, addr string, handler interface{}) (*Server, error) {
	isOpen := isPortOpen(addr)
	if isOpen {
		return nil, fmt.Errorf("port is open:%v", addr)
	}
	rpcServ := rpc.NewServer()
	err := rpcServ.RegisterName(name, handler)
	if err != nil {
		logger.Error("register name error:%v %v", name, err)
		return nil, err
	}
	rpcHandler := rpcServ.WebsocketHandler([]string{"*"})
	rpcServ.SetBatchLimits(BatchRequestLimit, BatchResponseMaxSize)
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        rpcHandler,
		ReadTimeout:    HttpReadTimeOut,
		WriteTimeout:   HttpWriteTimeOut,
		MaxHeaderBytes: MaxHeaderBytes,
		IdleTimeout:    3 * time.Hour,
	}
	return &Server{httpServer: httpServer, name: name}, nil
}

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		logger.Info("rpc server exit now: %v", err)
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	if s.httpServer != nil {
		err := s.httpServer.Shutdown(context.TODO())
		if err != nil {
			logger.Error("rpc server shutdown %v error:%v", s.name, err)
			return err
		}
	}
	return nil
}

func isPortOpen(endpoint string) bool {
	split := strings.Split(endpoint, ":")
	if len(split) != 2 {
		return true
	}
	listener, err := net.Listen("tcp", ":"+split[1])
	if err != nil {
		return true
	}
	_ = listener.Close()
	return false
}
