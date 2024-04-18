package service

import (
	"fmt"
	"github.com/NikitosnikN/balance-api/internal/adapters/nodepool"
	"github.com/NikitosnikN/balance-api/internal/app"
	"github.com/NikitosnikN/balance-api/internal/app/query"
	"github.com/NikitosnikN/balance-api/internal/config"
	"github.com/NikitosnikN/balance-api/internal/ports/rest"
	"log"
	"net/http"
)

func NewApplication(pool *nodepool.NodePool) *app.Application {
	return &app.Application{
		Queries: app.Queries{
			FetchBalance: query.NewFetchBalanceHandler(pool.EthGetBalance),
			IsPoolAlive:  query.NewIsPoolAliveHandler(pool.IsAlive),
		},
	}
}

func RunApplication(cfg *config.Config) error {
	log.Println("Starting application")
	nodeList := nodepool.NewLinkedList()

	for _, node := range cfg.Rpcs {
		nodeList.Insert(nodepool.NewNode(node.Name, node.Url))
	}

	pool, err := nodepool.NewNodePool(nodeList, cfg.WorkerInterval)

	if err != nil {
		return err
	}

	application := NewApplication(pool)
	restHandler := rest.Handler(application)
	log.Println("Running HTTP server on port", cfg.HTTPPort)
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.HTTPPort), restHandler)
}
