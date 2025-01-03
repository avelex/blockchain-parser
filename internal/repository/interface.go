package repository

import (
	"context"

	"github.com/avelex/blockchain-parser/internal/types"
)

type Repository interface {
	GetTransactions(ctx context.Context, address string) ([]types.Transaction, error)
	SaveTransactions(ctx context.Context, address string, transactions []types.Transaction) error
}
