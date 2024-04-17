package nodepool

type Node struct {
	Name        string
	Url         string
	IsAlive     bool
	LatestBlock uint
}
