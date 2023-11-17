package p2p

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/libp2p/go-libp2p"
	kaddht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/muxer/yamux"
	tls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"github.com/libp2p/go-libp2p/p2p/transport/tcp"
	"github.com/libp2p/go-libp2p/p2p/transport/websocket"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

type LibP2p struct {
	discovery  *Discovery
	ctx        context.Context
	Messages   chan *Msg
	host       host.Host
	peerId     peer.ID
	handler    *Handler
	cancel     context.CancelFunc
	topic      *pubsub.Topic
	sub        *pubsub.Subscription
	Listener   *Listener
	cfg        *P2pConfig
	bootstraps []string
}

func NewLibP2p(cfg *P2pConfig) (*LibP2p, error) {
	if cfg == nil {
		cfg = NewP2pConfig("", 4001, nil)
	}
	var privateKey crypto.PrivKey
	var err error
	if cfg.PrivateKey == "" {
		privateKey, _, err = crypto.GenerateECDSAKeyPair(rand.Reader)
		if err != nil {
			logger.Error("generate ed25519 private key error:%v", err)
			return nil, err
		}
	} else {
		privateBytes, err := hex.DecodeString(cfg.PrivateKey)
		if err != nil {
			logger.Error("decode ed25519 private key error:%v", err)
			return nil, err
		}
		privateKey, err = crypto.UnmarshalSecp256k1PrivateKey(privateBytes)
		if err != nil {
			logger.Error("unmarshal ed25519 private key error:%v", err)
			return nil, err
		}
	}
	muxers := libp2p.Muxer(zkbtcMuxerName, yamux.DefaultTransport)
	security := libp2p.Security(tls.ID, tls.New)
	listenAddres := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%v", cfg.Port),
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%v/ws", cfg.Port),
	)
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
		libp2p.Transport(websocket.New),
	)
	var dht *kaddht.IpfsDHT
	newDHT := func(h host.Host) (routing.PeerRouting, error) {
		var err error
		dht, err = kaddht.New(context.Background(), h)
		return dht, err
	}
	p2pRouting := libp2p.Routing(newDHT)
	node, err := libp2p.New(
		libp2p.Identity(privateKey),
		muxers,
		transports,
		security,
		listenAddres,
		p2pRouting,
	)
	if err != nil {
		logger.Error("create libp2p error:%v", err)
		return nil, err
	}
	discovery, err := NewDiscovery(node, nil)
	if err != nil {
		logger.Error("setup Discovery error:%v", err)
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	listener := NewListener()
	node.Network().Notify(listener)
	ps, err := pubsub.NewGossipSub(ctx, node)
	if err != nil {
		logger.Error("create gossipsub error:%v", err)
		cancel()
		return nil, err
	}
	topic, err := ps.Join(zkbtcLibP2pProtocol)
	if err != nil {
		logger.Error("create gossipsub topic error:%v", err)
		cancel()
		return nil, err
	}
	sub, err := topic.Subscribe()
	if err != nil {
		logger.Error("create gossipsub sub error:%v", err)
		cancel()
		return nil, err
	}
	err = dht.Bootstrap(ctx)
	if err != nil {
		logger.Error("dht bootstraps error:%v", err)
		cancel()
		return nil, err
	}
	logger.Debug("new libp2p success: %v", node.ID())
	return &LibP2p{
		discovery:  discovery,
		host:       node,
		ctx:        ctx,
		cfg:        cfg,
		cancel:     cancel,
		topic:      topic,
		sub:        sub,
		peerId:     node.ID(),
		bootstraps: cfg.Bootstrap,
		Listener:   listener,
		Messages:   make(chan *Msg, 1024),
	}, nil
}
func (p *LibP2p) SayHello(addr string) error {
	timestamp := time.Now().Unix()
	err := p.Broadcast(&Msg{
		Timestamp: &timestamp,
		Hello: &Hello{
			Address: &addr,
		}})
	if err != nil {
		return err
	}
	return err
}

func (p *LibP2p) AddPeer(addr string) error {
	// /ip4/127.0.0.1/tcp/4001/p2p/16Uiu2HAmGn4YKZcptkyue47yidbnKpm9wYqyaoPmeeCRUHCfJfZJ
	peerAddrInfo, err := peer.AddrInfoFromString(addr)
	if err != nil {
		logger.Error("create peer info error:%v", err)
		return err
	}
	err = p.host.Connect(p.ctx, *peerAddrInfo)
	if err != nil {
		logger.Warn("connect peer error:%v", err)
		return err
	}
	logger.Debug("connected peer : %v", addr)
	return nil
}

func (p *LibP2p) Broadcast(msg *Msg) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = p.topic.Publish(p.ctx, msgBytes)
	return nil
}

func (p *LibP2p) bootstrap(peers []string) {
	var waitGroup sync.WaitGroup
	for _, addr := range peers {
		waitGroup.Add(1)
		go func(addr string) {
			defer waitGroup.Done()
			err := p.AddPeer(addr)
			if err != nil {
				logger.Warn("connect addr %v error:%v,retry again", addr, err)
				peerAddr, err := multiaddr.NewMultiaddr(addr)
				if err != nil {
					logger.Error("parse addr %v error:%v", addr, err)
					return
				}
				muaAddr, peerId := peer.SplitAddr(peerAddr)
				p.Listener.AddDisConn(peerId, muaAddr)
			}
		}(addr)
	}
	waitGroup.Wait()
}

func (p *LibP2p) Run() {
	go p.bootstrap(p.cfg.Bootstrap)
	go p.readMessage()
	go p.retryConn()
}

func (p *LibP2p) retryConn() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.Listener.disConnPeers.Range(func(key, value interface{}) bool {
				peerId := key.(peer.ID)
				addr := value.(multiaddr.Multiaddr)
				if p.Listener.CanConn() {
					err := p.AddPeer(fmt.Sprintf("%v/p2p/%v", addr.String(), peerId.String()))
					if err != nil {
						//logger.Error("connect peerAddr %v %v error:%v", peerId, addr.String(), err)
					}
				}
				return true
			})
		}
	}
}

func (p *LibP2p) MsgChan() <-chan *Msg {
	return p.Messages
}

func (p *LibP2p) Close() error {
	p.cancel()
	err := p.topic.Close()
	if err != nil {
	}
	p.sub.Cancel()
	err = p.discovery.Close()
	if err != nil {
	}
	err = p.host.Close()
	if err != nil {
	}
	return nil
}

func (p *LibP2p) readMessage() {
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			msg, err := p.sub.Next(p.ctx)
			if err != nil {
				logger.Error("read message error:%v", err)
				continue
			}
			if msg.GetFrom() == p.peerId {
				continue
			}
			message := &Msg{}
			err = proto.Unmarshal(msg.Data, message)
			if err != nil {
				logger.Error("unmarshal message error:%v", err)
				continue
			}
			logger.Debug("receive %v message: %v", msg.GetFrom(), message.String())
			p.Messages <- message

		}
	}
}

type P2pConfig struct {
	PrivateKey string
	Port       int
	Bootstrap  []string
}

func NewP2pConfig(privateKey string, port int, bootstrap []string) *P2pConfig {
	return &P2pConfig{
		PrivateKey: privateKey,
		Port:       port,
		Bootstrap:  bootstrap,
	}
}

func toBootstrap(bootstraps []string) ([]peer.AddrInfo, error) {
	var bootstrap []peer.AddrInfo
	for _, addr := range bootstraps {
		peerInfo, err := peer.AddrInfoFromString(addr)
		if err != nil {
			logger.Error("create peer info error:%v", err)
			return nil, err
		}
		bootstrap = append(bootstrap, *peerInfo)
	}
	return bootstrap, nil
}
