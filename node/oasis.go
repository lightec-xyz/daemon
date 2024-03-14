package node

import "github.com/lightec-xyz/daemon/store"

var _ IAgent = (*OasisAgent)(nil)

type OasisAgent struct {
	store store.IStore
	name  string
}

func NewOasisAgent() (*OasisAgent, error) {
	return &OasisAgent{
		name: "oasis",
	}, nil
}
func (o *OasisAgent) ScanBlock() error {
	//TODO implement me
	panic("implement me")
}

func (o *OasisAgent) ProofResponse(resp ZkProofResponse) error {
	//TODO implement me
	panic("implement me")
}

func (o *OasisAgent) Init() error {
	//TODO implement me
	panic("implement me")
}

func (o *OasisAgent) Close() error {
	//TODO implement me
	panic("implement me")
}

func (o *OasisAgent) Name() string {
	//TODO implement me
	panic("implement me")
}
