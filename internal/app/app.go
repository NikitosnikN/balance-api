package app

import "github.com/NikitosnikN/balance-api/internal/app/query"

type Queries struct {
	FetchBalance query.FetchBalanceHandler
}

type Application struct {
	Queries Queries
}
