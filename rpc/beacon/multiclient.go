package beacon

import (
	"fmt"
	"github.com/lightec-xyz/daemon/rpc/beacon/types"
	"github.com/prysmaticlabs/prysm/v5/api/server/structs"
	"sync"
)

type IMultiBeacon interface {
	IBeacon
	Next()
}

type MultiClient struct {
	*Client
	clients []*Client
	index   int
	lock    sync.Mutex
}

func (m *MultiClient) Next() {
	defer m.lock.Unlock()
	m.lock.Lock()
	m.index++
	if m.index > len(m.clients)-1 {
		m.index = 0
	}
	m.Client = m.clients[m.index]

}
func (m *MultiClient) GetBlindedBlock(slot uint64) (types.BindBlockResp, error) {
	return m.Client.GetBlindedBlock(slot)
}

func (m *MultiClient) Eth1MapToEth2(slot uint64) (*Eth1MapToEth2, error) {
	return m.Client.Eth1MapToEth2(slot)
}

func (m *MultiClient) Bootstrap(slot uint64) (*types.BootstrapResp, error) {
	return m.Client.Bootstrap(slot)
}

func (m *MultiClient) BootstrapByRoot(root string) (*types.BootstrapResp, error) {
	return m.BootstrapByRoot(root)
}

func (m *MultiClient) BeaconHeaders(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	return m.Client.BeaconHeaders(slot)
}

func (m *MultiClient) FinalizedSlot() (uint64, error) {
	return m.Client.FinalizedSlot()
}

func (m *MultiClient) FinalizedPeriod() (uint64, error) {
	return m.Client.FinalizedPeriod()
}

func (m *MultiClient) LightClientUpdates(period, count uint64) ([]types.LightClientUpdateResp, error) {
	return m.Client.LightClientUpdates(period, count)
}

func (m *MultiClient) BeaconHeaderBySlot(slot uint64) (*structs.GetBlockHeaderResponse, error) {
	return m.Client.BeaconHeaderBySlot(slot)
}

func (m *MultiClient) BeaconHeaderByRoot(root string) (*structs.GetBlockHeaderResponse, error) {
	return m.Client.BeaconHeaderByRoot(root)
}

func (m *MultiClient) GetFinalityUpdate() (types.LightClientFinalityUpdateResp, error) {
	return m.Client.GetFinalityUpdate()
}

func (m *MultiClient) RetrieveBeaconHeaders(start, end uint64) ([]*structs.BeaconBlockHeader, error) {
	return m.Client.RetrieveBeaconHeaders(start, end)
}

func NewMultiClient(url ...string) (IMultiBeacon, error) {
	if len(url) == 0 {
		return nil, fmt.Errorf("no url")
	}
	var clients []*Client
	for _, u := range url {
		client, err := NewClient(u)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return &MultiClient{
		clients: clients,
		Client:  clients[0],
		index:   0,
	}, nil
}
