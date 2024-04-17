package query

import (
	"context"
)

type (
	FetchBalanceQuery struct {
		Address  string
		BlockTag string
	}
	FetchBalanceResult struct {
		Balance string
	}
	FetchBalanceHandler struct {
	}
)

func NewFetchBalanceHandler() *FetchBalanceHandler {
	return &FetchBalanceHandler{}
}

func (h *FetchBalanceHandler) Handle(ctx context.Context, query FetchBalanceQuery) (*FetchBalanceResult, error) {
	return &FetchBalanceResult{}, nil

}
