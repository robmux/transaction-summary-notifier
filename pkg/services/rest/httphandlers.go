package rest

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/summary"
	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
)

type (
	TransactionsManager interface {
		LoadTransactions(ctx context.Context, userID uint64) ([]transactions.TransactionDetail, error)
	}

	SummaryService interface {
		GetSummary(ctx context.Context) (*summary.GeneralSummary, error)

		GetTotalBalanceInAccount(ctx context.Context, transactions []transactions.TransactionDetail) decimal.Decimal

		GetNumberOfTransactionsGroupedByMonth(ctx context.Context, transactions []transactions.TransactionDetail) map[uint8]summary.TransactionsByMonth
		GetAverageCreditAndDebit(ctx context.Context, transactions []transactions.TransactionDetail) summary.AveragesByMonth

		GetAverageDebit(ctx context.Context, transactions []transactions.TransactionDetail) summary.AmountDetail
		GetAverageCredit(ctx context.Context, transactions []transactions.TransactionDetail) summary.AmountDetail
	}
)

type Handler struct {
	TransactionsSrv TransactionsManager
	SummarySrv      SummaryService
}

func New(txManager TransactionsManager, summarySrv SummaryService) *Handler {
	return &Handler{
		TransactionsSrv: txManager,
		SummarySrv:      summarySrv,
	}
}

func (h *Handler) Ping(ctx *gin.Context) error {
	ctx.String(200, "pong")
	return nil
}

type appHandler func(ctx *gin.Context) error

// makeHTTPHandlers converts an handler with err to a gin handler
func makeHTTPHandler(h appHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h(ctx)
		if err != nil {
			ctx.String(500, err.Error())
		}

		// Nothing here, because h should have already sent the response
	}
}
