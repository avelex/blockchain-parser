package ethclient

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/avelex/blockchain-parser/internal/jsonrpc"
)

const (
	blockNumberMethod        = "eth_blockNumber"
	blockByNumberMethod      = "eth_getBlockByNumber"
	transactionReceiptMethod = "eth_getTransactionReceipt"
)

// NOTE: add rate-limit
type Client struct {
	id  string
	url string
	rpc *jsonrpc.Client
}

func New(url string) *Client {
	return &Client{
		id:  randomID(),
		url: url,
		rpc: jsonrpc.NewClient(),
	}
}

func (c *Client) BlockNumber(ctx context.Context) (int, error) {
	req := jsonrpc.NewEmptyRequest(blockNumberMethod, c.id)

	resp, err := c.rpc.Call(ctx, c.url, req)
	if err != nil {
		return 0, fmt.Errorf("failed to call %s: %w", blockNumberMethod, err)
	}

	blockHex, ok := resp.Result.(string)
	if !ok {
		return 0, fmt.Errorf("failed to parse %s response", blockNumberMethod)
	}

	blockInt, err := parseHexInt(blockHex)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s block number", blockNumberMethod)
	}

	return int(blockInt), nil
}

func (c *Client) BlockHeaderByNumber(ctx context.Context, number int) (*BlockHeader, error) {
	numberHex := "0x" + strconv.FormatInt(int64(number), 16)
	params := []any{
		numberHex,
		false, // don't include transaction objects, only hashes
	}

	req := jsonrpc.NewRequest(blockByNumberMethod, params, c.id)

	resp, err := c.rpc.Call(ctx, c.url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call %s: %w", blockByNumberMethod, err)
	}

	header, err := blockHeaderFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", blockByNumberMethod, err)
	}

	return header, nil
}

func (c *Client) TransactionReceipt(ctx context.Context, hash string) (*TransactionReceipt, error) {
	req := jsonrpc.NewRequest(transactionReceiptMethod, []any{hash}, c.id)

	resp, err := c.rpc.Call(ctx, c.url, req)
	if err != nil {
		return nil, fmt.Errorf("failed to call %s: %w", transactionReceiptMethod, err)
	}

	receipt, err := transactionReceiptFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s response: %w", transactionReceiptMethod, err)
	}

	return receipt, nil
}

func randomID() string {
	buff := make([]byte, 4)
	if _, err := rand.Read(buff); err != nil {
		return "1"
	}
	return new(big.Int).SetBytes(buff).String()
}
