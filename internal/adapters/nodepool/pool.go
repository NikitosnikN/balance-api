package nodepool

import (
	"context"
	"errors"
	"github.com/NikitosnikN/balance-api/internal/common"
	"github.com/NikitosnikN/balance-api/pkg/linked_list"
	"log"
	"time"
)

type NodePool struct {
	nodes          *linked_list.LinkedList[Node]
	workerInterval time.Duration
}

func NewNodePool(nodes *linked_list.LinkedList[Node], workerInterval time.Duration) (*NodePool, error) {
	if len(nodes.GetElements()) == 0 {
		return nil, errors.New("nodes is empty")
	}
	pool := NodePool{
		nodes:          nodes,
		workerInterval: workerInterval,
	}

	err := pool.runWorker()

	if err != nil {
		return nil, err
	}

	return &pool, nil
}

func updateNodeState(node *Node) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	blockNumber, err := node.EthGetBlockNumber(ctx)

	if err != nil {
		node.IsAlive = false
		return
	}

	number, err := common.HexToDecimal(blockNumber)

	if err != nil {
		node.IsAlive = false
		return
	}

	node.LatestBlock = uint(number)
	node.IsAlive = true
}

func (n *NodePool) runIteration() {
	nodes := n.nodes.GetElements()

	if len(nodes) == 0 {
		return
	}

	for _, node := range nodes {
		go updateNodeState(&node)
	}
}

func (n *NodePool) runWorker() error {
	log.Println("Starting node pool worker")
	ticker := time.NewTicker(n.workerInterval)

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
	el := n.nodes.Next()

	if el == nil {
		return "", errors.New("cannot get node")
	}

	node := el.Data()

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
