package p2p

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lightec-xyz/daemon/logger"
	"github.com/multiformats/go-multiaddr"
	"sync"
)

type Listener struct {
	disConnPeers *sync.Map
	connNum      int
}

func (l *Listener) Listen(n network.Network, multiaddr multiaddr.Multiaddr) {
	logger.Debug("listen %v", multiaddr.String())
}

func (l *Listener) ListenClose(n network.Network, multiaddr multiaddr.Multiaddr) {
	logger.Debug("listen close %v", multiaddr.String())
}

func (l *Listener) Connected(n network.Network, conn network.Conn) {
	logger.Debug("connected %v %v", conn.RemotePeer(), conn.RemoteMultiaddr())
	if _, ok := l.disConnPeers.Load(conn.RemotePeer()); ok {
		l.disConnPeers.Delete(conn.RemotePeer())
	}
	l.connNum++
}

func (l *Listener) Disconnected(n network.Network, conn network.Conn) {
	logger.Debug("disconnected %v", conn.RemotePeer())
	if _, ok := l.disConnPeers.Load(conn.RemotePeer()); !ok {
		l.disConnPeers.Store(conn.RemotePeer(), conn.RemoteMultiaddr())
	}
	l.connNum--

}

func (l *Listener) AddDisConn(peer peer.ID, addr multiaddr.Multiaddr) {
	l.disConnPeers.Store(peer, addr)
}
func (l *Listener) CanConn() bool {
	return l.connNum < 100 // todo
}

func NewListener() *Listener {
	return &Listener{
		disConnPeers: new(sync.Map),
	}
}
