package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gorilla/websocket"
	"github.com/lightec-xyz/daemon/logger"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	WsConnPath = "/ws"
)

type Server struct {
	httpServer *http.Server
	name       string
}

func NewServer(name, addr string, handler interface{}, verify IVerify, wsHandler WsFn) (*Server, error) {
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
	var middlewareHandler http.Handler
	if wsHandler == nil {
		middlewareHandler = CORSHandler(VerifyHandler(rpcServer, verify))
	} else {
		middlewareHandler = WsConnHandler(VerifyHandler(rpcServer, verify), wsHandler)
	}
	rpcServer.SetBatchLimits(BatchRequestLimit, BatchResponseMaxSize)
	httpServer := &http.Server{
		Addr:           addr,
		Handler:        middlewareHandler,
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

func NewCustomWsServer(name, addr string, fn WsFn) (*Server, error) {
	isOpen := isPortOpen(addr)
	if isOpen {
		return nil, fmt.Errorf("port is open:%v", addr)
	}
	if fn == nil {
		return nil, fmt.Errorf("fn is nil")
	}
	rpcHandler := WsWrappHandler(fn)
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

type WsFn func(opt *WsOpt) error

type WsOpt struct {
	Id   string
	Conn *websocket.Conn
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

func WsWrappHandler(fn func(opt *WsOpt) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: wsHandshakeValidator([]string{"*"}),
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrade ws Conn error: %v", err)
			return
		}
		// todo
		err = fn(&WsOpt{Conn: conn})
		return
	})
}

func wsHandshakeValidator(allowedOrigins []string) func(*http.Request) bool {
	return func(r *http.Request) bool {
		for _, origin := range allowedOrigins {
			if origin == "*" {
				return true
			}
		}
		return false
	}
}
func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

func VerifyHandler(h http.Handler, verify IVerify) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if verify != nil {
			err := verifyJwt(r, verify)
			if err != nil {
				logger.Error("verify jwt error: %v", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func WsConnHandler(h http.Handler, fn func(opt *WsOpt) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == WsConnPath {
			// todo
			var upgrader = websocket.Upgrader{
				CheckOrigin: wsHandshakeValidator([]string{"*"}),
			}
			id := r.URL.Query().Get("id")
			conn, err := upgrader.Upgrade(w, r, nil)
			logger.Debug("new ws Conn coming: %v id: %v", r.URL.Path, id)
			if err != nil {
				log.Printf("upgrade ws Conn error: %v", err)
				return
			}
			err = fn(&WsOpt{Id: id, Conn: conn})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
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

func pathPermission(method string) (Permission, error) {
	switch method {
	case RemoveProofPath:
		return JwtPermission, nil
	default:
		return NonePermission, nil
	}
}

// todo
func verifyJwt(request *http.Request, verify IVerify) error {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		logger.Error("read body error: %v", err)
		return err
	}
	method, err := getMethod(bodyBytes)
	if err != nil {
		return err
	}
	permission, err := pathPermission(method)
	if err != nil {
		return err
	}
	if permission != NonePermission {
		authorization := request.Header.Get(AuthorizationHeader)
		jwt, err := verify.VerifyJwt(authorization)
		if err != nil {
			return fmt.Errorf("invalid jwt token")
		}
		if permission != jwt.Permission {
			return fmt.Errorf("permission not match")
		}
		return nil
	}

	request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return nil
}

func getMethod(data []byte) (string, error) {
	var req ReqMsg
	err := json.Unmarshal(data, &req)
	if err != nil {
		return "", err
	}
	return req.Method, nil
}

type ReqMsg struct {
	Method  string `json:"method"`
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
}
