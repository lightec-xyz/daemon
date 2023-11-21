package agent

type IAgent interface {
	ScanBlock(height int64) (int64, error)
	Init() error
}
