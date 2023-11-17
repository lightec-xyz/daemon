package p2p

import (
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/lightec-xyz/daemon/logger"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) StreamHandler(s network.Stream) {
	if s.Protocol() != zkbtcLibP2pProtocol {
		logger.Warn("find unSupport p2p protocol:%v", s.Protocol())
		return
	}

}

func (h *Handler) Close() error {
	return nil
}
