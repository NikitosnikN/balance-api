package service

import (
	"github.com/NikitosnikN/balance-api/internal/app"
	"github.com/NikitosnikN/balance-api/internal/app/query"
	"github.com/NikitosnikN/balance-api/internal/config"
)

func NewApplication() *app.Application {
	return &app.Application{
		Queries: app.Queries{
			FetchBalance: query.FetchBalanceHandler{},
		},
	}
}

func RunApplication(cfg *config.Config) error {
	return nil
}
