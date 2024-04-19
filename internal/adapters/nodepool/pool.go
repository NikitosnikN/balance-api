package nodepool

import (
	"context"
	"errors"
	"github.com/NikitosnikN/balance-api/internal/common"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

type Metrics interface {
	NodePoolBlockHeight() *prometheus.GaugeVec
	NodePoolLiveness() *prometheus.GaugeVec
}

type NodePool struct {
	nodes          *LinkedList
	workerInterval time.Duration
	metrics        Metrics
}

func NewNodePool(nodes *LinkedList, workerInterval time.Duration, metrics Metrics) (*NodePool, error) {
	if len(nodes.GetElements()) == 0 {
		return nil, errors.New("nodes is empty")
	}
	pool := NodePool{
		nodes:          nodes,
		workerInterval: workerInterval,
		metrics:        metrics,
	}

	err := pool.runWorker()

	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func updateNodeState(node *Node, metrics Metrics) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	blockNumber, err := node.EthGetBlockNumber(ctx)

	if err != nil {
		log.Printf("Node %s is not alive.\n", node.Name)
		node.IsAlive = false
		metrics.NodePoolLiveness().WithLabelValues(node.Name).Set(0)
		return
	}

	number, err := common.HexToDecimal(blockNumber)

	if err != nil {
		log.Printf("Node %s is not alive.\n", node.Name)
		node.IsAlive = false
		metrics.NodePoolLiveness().WithLabelValues(node.Name).Set(0)
		return
	}

	node.LatestBlock = uint(number)
	node.IsAlive = true

	metrics.NodePoolBlockHeight().WithLabelValues(node.Name).Set(float64(number))
	metrics.NodePoolLiveness().WithLabelValues(node.Name).Set(1)
}

func (n *NodePool) runIteration() {
	nodes := n.nodes.GetElements()

	if len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		go updateNodeState(node, n.metrics)
	}
}

func (n *NodePool) runWorker() error {
	log.Println("Starting node pool worker")
	ticker := time.NewTicker(n.workerInterval)

	// warming nodes' status
	go n.runIteration()

	if len(n.nodes.GetElements()) == 0 {
		return errors.New("node pool is empty")
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				go n.runIteration()
			}
		}
	}()

	return nil
}

func (n *NodePool) EthGetBalance(ctx context.Context, address string, blockTag string) (string, error) {
	node := n.nodes.NextAliveNode()

	if node == nil {
		return "", errors.New("cannot get node")
	}

	return node.EthGetBalance(ctx, address, blockTag)
}

func (n *NodePool) IsAlive() (bool, error) {
	nodes := n.nodes.GetElements()
	isAlive := false

	for _, node := range nodes {
		if node.IsAlive {
			isAlive = true
			break
		}
	}

	return isAlive, nil
}
