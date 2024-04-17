package query

import (
	"context"
	"errors"
	"github.com/NikitosnikN/balance-api/internal/common"
)

type (
	FetchBalanceQuery struct {
		Address  string `json:"address"`
		BlockTag string `json:"blockTag"`
	}
	FetchBalanceResult struct {
		Balance string `json:"balance"`
	}
	FetchBalanceHandler struct {
		getBalance func(ctx context.Context, address string, blockTag string) (string, error)
	}
)

func NewFetchBalanceHandler(
	getBalance func(ctx context.Context, address string, blockTag string) (string, error),
) FetchBalanceHandler {
	return FetchBalanceHandler{
		getBalance: getBalance,
	}
}

func (h *FetchBalanceHandler) Handle(ctx context.Context, query FetchBalanceQuery) (*FetchBalanceResult, error) {
	ethAddressValid := common.IsEthAddress(query.Address)

	if !ethAddressValid {
		return nil, errors.New("invalid eth address")
	}

	blockTagValid := common.IsBlockTag(query.BlockTag)

	if !blockTagValid {
		return nil, errors.New("invalid block tag")
	}

	balance, err := h.getBalance(ctx, query.Address, query.BlockTag)

	if err != nil {
		return nil, err
	}

	return &FetchBalanceResult{
		Balance: balance,
	}, nil

}
