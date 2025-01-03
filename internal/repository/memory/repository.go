package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/avelex/blockchain-parser/internal/types"
)

type Repository struct {
	mu          *sync.RWMutex
	subscribers map[string][]types.Transaction
}

func New() *Repository {
	return &Repository{
		mu:          &sync.RWMutex{},
		subscribers: make(map[string][]types.Transaction),
	}
}

func (r *Repository) GetTransactions(ctx context.Context, address string) ([]types.Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tx, ok := r.subscribers[address]
	if !ok {
		return nil, fmt.Errorf("address not found")
	}

	return tx, nil
}

func (r *Repository) SaveTransactions(ctx context.Context, address string, transactions []types.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscribers[address] = append(r.subscribers[address], transactions...)

	return nil
}
