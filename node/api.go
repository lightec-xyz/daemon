package node

type API interface {
	Version() (DaemonInfo, error)
}
