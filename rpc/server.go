package rpc

import (
	"context"
	"github.com/ethereum/go-ethereum/rpc"
	"log"
	"net/http"
)

type Server struct {
	rpcServer *rpc.Server
	server    *http.Server
}

func NewServer(addr string) (*Server, error) {
	rpcServ := rpc.NewServer()
	err := rpcServ.RegisterName("zkbtc", new(Handler))
	if err != nil {
		log.Fatal(err)
	}
	//todo
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

		}
	}
	return nil
}
