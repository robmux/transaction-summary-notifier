package transactions

import (
	"context"

	"github.com/shopspring/decimal"
)

type TransactionDetail struct {
	ID   uint64
	Date MonthDay

	TransactionAmount decimal.Decimal
}

type MonthDay struct {
	Month uint8
	Day   uint8
}

type TransactionSrv struct {
	transactionsGetter TransactionsGetter
}

func New(transactionsGetter TransactionsGetter) *TransactionSrv {
	return &TransactionSrv{
		transactionsGetter: transactionsGetter,
	}
}

func (s *TransactionSrv) LoadTransactions(ctx context.Context, userID uint64) ([]TransactionDetail, error) {
	return s.transactionsGetter.GetUserTransactions(ctx, userID)
}
