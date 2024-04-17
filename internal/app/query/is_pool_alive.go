package query

import (
	"context"
)

type (
	IsPoolAliveQuery struct {
	}
	IsPoolAliveResult struct {
		Alive bool
	}
	IsPoolAliveHandler struct {
		isAlive func() (bool, error)
	}
)

func NewIsPoolAliveHandler(
	isAlive func() (bool, error),
) IsPoolAliveHandler {
	return IsPoolAliveHandler{
		isAlive: isAlive,
	}
}

func (h *IsPoolAliveHandler) Handle(ctx context.Context, query IsPoolAliveQuery) (*IsPoolAliveResult, error) {
	alive, err := h.isAlive()

	if err != nil {
		return nil, err
	}

	return &IsPoolAliveResult{
		Alive: alive,
	}, nil

}
