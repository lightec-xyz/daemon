package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/lightec-xyz/daemon/rpc/ws"
	"strings"
)

type CustomServer struct {
	service *Service
	conn    *ws.Conn
}

func NewCustomServer(url string, handler interface{}) (*CustomServer, error) {
	service := NewService(handler)
	wsConn, err := ws.NewClientConn(url, service.Call, nil, false)
	if err != nil {
		logger.Error("new client conn error:%v %v", url, err)
		return nil, err
	}
	return &CustomServer{
		service: service,
		conn:    wsConn,
	}, nil
}

func (s *CustomServer) Run() error {
	s.conn.Run()
	return nil
}

func (s *CustomServer) Close() error {
	if s.service != nil {
		s.service.Close()
	}
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}

type Service struct {
	rpcService *ws.Service
}

func NewService(handler interface{}) *Service {
	rpcService := ws.NewService(handler)
	return &Service{
		rpcService: rpcService,
	}
}
func (s *Service) Close() error {
	return nil
}

func (s *Service) Call(req ws.Message) (ws.Message, error) {
	exists := s.CheckMethod(req.Method)
	if !exists {
		logger.Error("no such method: %s", req.Method)
		return ws.NewErrorMsg(req.Id, req.Method, fmt.Sprintf("no such method %s", req.Method)), fmt.Errorf("no such method: %s", req.Method)
	}

	result, err := s.rpcService.Call(req.Method, req.Data)
	if err != nil {
		logger.Error("call error:%v %v", req.Method, err)
		return ws.NewErrorMsg(req.Id, req.Method, err.Error()), err
	} else {
		if result == nil {
			return ws.NewRespMessage(req.Id, req.Method, nil), nil
		}
		resultBytes, err := json.Marshal(result)
		if err != nil {
			logger.Error("marshal error:%v %v", req.Method, err)
			return ws.NewErrorMsg(req.Id, req.Method, err.Error()), err
		}
		return ws.NewRespMessage(req.Id, req.Method, resultBytes), nil
	}
}

func (s *Service) CheckMethod(method string) bool {
	_, after, _ := strings.Cut(method, "_")
	exists := s.rpcService.Check(after)
	return exists
}
