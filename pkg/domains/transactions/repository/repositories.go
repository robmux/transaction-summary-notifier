package repository

import (
	"context"
	
	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
)

type TransactionsGetter interface {
	GetUserTransactions(ctx context.Context, userID uint64) ([]transactions.TransactionDetail, error)
}
