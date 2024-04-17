package app

import "github.com/NikitosnikN/balance-api/internal/app/query"

type Queries struct {
	FetchBalance query.FetchBalanceHandler
	IsPoolAlive  query.IsPoolAliveHandler
}

type Application struct {
	Queries Queries
}
