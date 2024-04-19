package nodepool

import (
	"context"
	"errors"
	"github.com/NikitosnikN/balance-api/internal/common"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"time"
)

type Metrics interface {
	NodePoolBlockHeight() *prometheus.GaugeVec
	NodePoolLiveness() *prometheus.GaugeVec
}

type NodePool struct {
	nodes          []*Node
	workerInterval time.Duration
	metrics        Metrics
}

func NewNodePool(nodes []*Node, workerInterval time.Duration, metrics Metrics) (*NodePool, error) {
	if len(nodes) == 0 {
		return nil, errors.New("node list is empty")
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
	if len(n.nodes) == 0 {
		return
	}

	for _, node := range n.nodes {
		go updateNodeState(node, n.metrics)
	}
}

func (n *NodePool) runWorker() error {
	log.Println("Starting node pool worker")
	ticker := time.NewTicker(n.workerInterval)

	// warming nodes' status
	go n.runIteration()

	if len(n.nodes) == 0 {
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

func (n *NodePool) GetAliveNodes() []*Node {
	var nodes []*Node

	for i, node := range n.nodes {
		if node.IsAlive {
			nodes = append(nodes, n.nodes[i])
		}
	}

	return nodes
}

func (n *NodePool) EthGetBalance(ctx context.Context, address string, blockTag string) (string, error) {
	nodes := n.GetAliveNodes()

	if len(nodes) == 0 {
		return "", errors.New("no alive nodes in node pool")
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	wg := sync.WaitGroup{}
	resultsCh := make(chan string)
	results := map[string]int{}

	fetchFunc := func(wg *sync.WaitGroup, node *Node, resultsCh chan string) {
		defer wg.Done()
		result, err := node.EthGetBalance(ctx, address, blockTag)
		if err == nil && result != "" {
			resultsCh <- result
		}
	}

	// listen channel, update results
	go func() {
		for {
			select {
			case result := <-resultsCh:
				results[result]++
			case <-ctx.Done():
				return
			}
		}
	}()

	// fire tasks
	for _, node := range nodes {
		wg.Add(1)
		go fetchFunc(&wg, node, resultsCh)
	}

	// wait for execution, close channel
	wg.Wait()
	close(resultsCh)

	// get true balance from results
	result, err := determineTrueBalance(results)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (n *NodePool) IsAlive() (bool, error) {
	return len(n.GetAliveNodes()) > 0, nil
}
