package parser

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/avelex/blockchain-parser/config"
	"github.com/avelex/blockchain-parser/internal/ethclient"
	"github.com/avelex/blockchain-parser/internal/repository"
	"github.com/avelex/blockchain-parser/internal/types"
)

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(ctx context.Context, address string) []types.Transaction
}

type BlockchainParser struct {
	subMu      *sync.RWMutex
	subscriber map[string]struct{}

	currentBlock atomic.Int64

	conf   config.Config
	client *ethclient.Client
	repo   repository.Repository
}

func New(conf config.Config, client *ethclient.Client, repo repository.Repository) *BlockchainParser {
	return &BlockchainParser{
		subMu:        &sync.RWMutex{},
		subscriber:   make(map[string]struct{}),
		currentBlock: atomic.Int64{},
		conf:         conf,
		client:       client,
		repo:         repo,
	}
}

func (p *BlockchainParser) GetCurrentBlock() int {
	return int(p.currentBlock.Load())
}

func (p *BlockchainParser) Subscribe(address string) bool {
	address = strings.ToLower(address)

	p.subMu.Lock()
	defer p.subMu.Unlock()

	if _, ok := p.subscriber[address]; ok {
		return false
	}

	p.subscriber[address] = struct{}{}

	return true
}

func (p *BlockchainParser) GetTransactions(ctx context.Context, address string) []types.Transaction {
	address = strings.ToLower(address)

	tx, err := p.repo.GetTransactions(ctx, address)
	if err != nil {
		return []types.Transaction{}
	}
	return tx
}

func (p *BlockchainParser) Start(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	// NOTE: possible for parallel blocks processing
	sub := make(chan int)

	// start listen new blocks
	go func() {
		defer wg.Done()
		p.listenBlocks(ctx, sub)
	}()

	// process blocks for transactions
	go func() {
		defer wg.Done()
		p.processBlocks(ctx, sub)
	}()

	<-ctx.Done()

	wg.Wait()
	close(sub)

	return nil
}

func (p *BlockchainParser) listenBlocks(ctx context.Context, pub chan<- int) {
	ticker := time.NewTicker(p.conf.BlocksInterval)
	defer ticker.Stop()

	var startBlock int

	if p.conf.StartBlock != 0 {
		slog.Info("Start from block", "number", p.conf.StartBlock)
		startBlock = p.conf.StartBlock
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			currentBlock, err := p.client.BlockNumber(ctx)
			if err != nil {
				slog.Warn("failed to get current block number", "error", err)
				continue
			}

			if startBlock == 0 {
				startBlock = currentBlock
			} else if startBlock == currentBlock {
				slog.Info("No new blocks, wait for next block...", "startBlock", startBlock, "currentBlock", currentBlock)
				continue
			}

			for i := startBlock; i < currentBlock; i++ {
				pub <- i
			}

			startBlock = currentBlock
		}
	}
}

// NOTE: add processing for failed blocks and transactions
func (p *BlockchainParser) processBlocks(ctx context.Context, sub <-chan int) {
	for blockNumber := range sub {
		if isContextDone(ctx) {
			slog.Info("Context done, stop processing blocks")
			return
		}

		start := time.Now()
		slog.Info("Processing block", "number", blockNumber)

		bh, err := p.client.BlockHeaderByNumber(ctx, blockNumber)
		if err != nil {
			slog.Error("failed to get block header", "number", blockNumber, "error", err)
			continue
		}

		txChan := make(chan string, 2)
		result := make(chan *ethclient.TransactionReceipt, 1)

		go func() {
			for i := 0; i < cap(txChan); i++ {
				go p.processTransactions(ctx, txChan, result)
			}
		}()

		go func() {
			for _, tx := range bh.Transactions {
				txChan <- tx
			}
		}()

		subTx := make(map[string][]types.Transaction)

		for i := 0; i < len(bh.Transactions); i++ {
			receipt := <-result

			if receipt.IsFailed() {
				continue
			}

			tx := types.NewTransaction(receipt.Hash, receipt.From, receipt.To, bh.Timestamp)

			if p.subscriberExists(receipt.From) {
				subTx[receipt.From] = append(subTx[receipt.From], tx)
			} else if p.subscriberExists(receipt.To) {
				subTx[receipt.To] = append(subTx[receipt.To], tx)
			}
		}

		// safe close, all transactions processed
		close(txChan)
		close(result)

		for address, txs := range subTx {
			if err := p.repo.SaveTransactions(ctx, address, txs); err != nil {
				slog.Error("failed to save transactions", "address", address, "error", err)
			}
		}

		p.currentBlock.Store(int64(blockNumber))

		slog.Info("Processed block", "number", blockNumber, "tx_count", len(bh.Transactions), "dur", time.Since(start))
	}
}

func (p *BlockchainParser) processTransactions(ctx context.Context, txChan <-chan string, result chan<- *ethclient.TransactionReceipt) {
	for txHash := range txChan {
		receipt, err := p.client.TransactionReceipt(ctx, txHash)
		if err != nil {
			result <- ethclient.EmptyFailedReceipt()
			continue
		}

		result <- receipt
	}
}

func (p *BlockchainParser) subscriberExists(address string) bool {
	p.subMu.RLock()
	defer p.subMu.RUnlock()

	_, ok := p.subscriber[address]

	return ok
}

func isContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
