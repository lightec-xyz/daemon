package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/lightec-xyz/daemon/logger"
)

type Discovery struct {
	h     host.Host
	close func() error
	peer  chan peer.AddrInfo
}

func NewDiscovery(h host.Host, peer chan peer.AddrInfo) (*Discovery, error) {
	DiscoveryServiceTag := "zkbtc-discovery"
	d := &Discovery{
		h:    h,
		peer: peer,
	}
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, d)
	err := s.Start()
	d.close = s.Close
	return d, err
}

func (d *Discovery) Close() error {
	return d.close()
}

func (d *Discovery) HandlePeerFound(pa peer.AddrInfo) {
	logger.Debug("discovered new peer %s", pa.ID)
	err := d.h.Connect(context.Background(), pa)
	if err == nil && d.peer != nil {
		d.peer <- pa
	}
}
