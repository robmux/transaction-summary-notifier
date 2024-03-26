package transactions

import (
	"context"
)

type TransactionsGetter interface {
	GetUserTransactions(ctx context.Context, userID uint64) ([]TransactionDetail, error)
}
