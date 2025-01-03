package ethclient

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/avelex/blockchain-parser/internal/jsonrpc"
)

// BlockHeader contains only used fields with transactions hashes
// for full block see https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_getblockbyhash
type BlockHeader struct {
	Number int `json:"number"`
	// 32 bytes hash
	Hash         string   `json:"hash"`
	Transactions []string `json:"transactions"`
	Timestamp    int64    `json:"timestamp"`
}

func blockHeaderFromResponse(r jsonrpc.Response) (*BlockHeader, error) {
	if r.Result == nil {
		return nil, fmt.Errorf("result is nil")
	}

	resultMap, ok := r.Result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("result is not a map")
	}

	var blockHeader BlockHeader

	if numberHex, ok := resultMap["number"]; ok {
		numberHex, ok := numberHex.(string)
		if !ok {
			return nil, fmt.Errorf("block number is not a string")
		}
		numberInt, err := parseHexInt(numberHex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse block number: %w", err)
		}
		blockHeader.Number = numberInt
	}

	if hash, ok := resultMap["hash"]; ok {
		blockHeader.Hash, ok = hash.(string)
		if !ok {
			return nil, fmt.Errorf("block hash is not a string")
		}
	}

	if transactions, ok := resultMap["transactions"]; ok {
		transactionsArr, ok := transactions.([]any)
		if !ok {
			return nil, fmt.Errorf("transactions is not an array")
		}

		for _, t := range transactionsArr {
			transactionHash, ok := t.(string)
			if !ok {
				return nil, fmt.Errorf("transaction hash is not a string")
			}
			blockHeader.Transactions = append(blockHeader.Transactions, transactionHash)
		}
	}

	if timestamp, ok := resultMap["timestamp"]; ok {
		timestampHex, ok := timestamp.(string)
		if !ok {
			return nil, fmt.Errorf("timestamp is not a string")
		}
		timestampInt, err := parseHexInt(timestampHex)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		blockHeader.Timestamp = int64(timestampInt)
	}

	return &blockHeader, nil
}

// TransactionReceipt contains only used fields
// for full receipt see https://ethereum.org/en/developers/docs/apis/json-rpc/#eth_gettransactionreceipt
type TransactionReceipt struct {
	// 1 (success) or 0 (failure)
	Status int `json:"status"`
	// 32 bytes hash
	Hash string `json:"transactionHash"`
	// 20 bytes address
	From string `json:"from"`
	// 20 bytes address
	To string `json:"to"`
}

func EmptyFailedReceipt() *TransactionReceipt {
	return &TransactionReceipt{
		Status: 0,
	}
}

func (t *TransactionReceipt) IsFailed() bool {
	return t.Status == 0
}

func transactionReceiptFromResponse(r jsonrpc.Response) (*TransactionReceipt, error) {
	if r.Result == nil {
		return nil, fmt.Errorf("result is nil")
	}

	resultMap, ok := r.Result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("result is not a map")
	}

	var receipt TransactionReceipt

	if status, ok := resultMap["status"]; ok {
		statusStr, ok := status.(string)
		if !ok {
			return nil, fmt.Errorf("status is not a string")
		}
		statusInt, err := parseHexInt(statusStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse status: %w", err)
		}
		receipt.Status = statusInt
	}

	if hash, ok := resultMap["transactionHash"]; ok {
		receipt.Hash, ok = hash.(string)
		if !ok {
			return nil, fmt.Errorf("transaction hash is not a string")
		}
	}

	if from, ok := resultMap["from"]; ok {
		receipt.From, ok = from.(string)
		if !ok {
			return nil, fmt.Errorf("from is not a string")
		}
	}

	if to, ok := resultMap["to"]; ok {
		receipt.To, ok = to.(string)
		if !ok {
			return nil, fmt.Errorf("to is not a string")
		}
	}

	return &receipt, nil
}

// parseHexInt parse hex string starting with 0x to int
func parseHexInt(s string) (int, error) {
	i, err := strconv.ParseInt(strings.TrimPrefix(s, "0x"), 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string: %w", err)
	}
	return int(i), nil
}
