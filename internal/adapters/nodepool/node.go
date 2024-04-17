package nodepool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type blockNumberResponse struct {
	Result string `json:"result"`
}

type getBalanceResponse struct {
	Result string `json:"result"`
}

type Node struct {
	Name        string
	Url         string
	IsAlive     bool
	LatestBlock uint
}

func NewNode(name string, url string) *Node {
	return &Node{
		Name:        name,
		Url:         url,
		IsAlive:     true,
		LatestBlock: 0,
	}
}

func (n *Node) EthGetBlockNumber(ctx context.Context) (string, error) {
	data := []byte(`{"jsonrpc": "2.0","method": "eth_blockNumber","params": [],"id": 1}`)

	req, err := http.NewRequestWithContext(ctx, "POST", n.Url, bytes.NewBuffer(data))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var response blockNumberResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		return "", err
	}

	return response.Result, nil
}

func (n *Node) EthGetBalance(ctx context.Context, address string, blockTag string) (string, error) {
	data := []byte(fmt.Sprintf(`{"jsonrpc": "2.0","method": "eth_getBalance","params": ["%s", "%s"],"id": 1}`, address, blockTag))

	req, err := http.NewRequestWithContext(ctx, "POST", n.Url, bytes.NewBuffer(data))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var response getBalanceResponse

	err = json.Unmarshal(body, &response)

	if err != nil {
		return "", err
	}

	return response.Result, nil
}
