package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
	"github.com/shopspring/decimal"
)

type (
	TransactionDetail struct {
		ID   uint64   `json:"id"`
		Date MonthDay `json:"date"`

		TransactionAmount decimal.Decimal `json:"transaction_amount"`
	}

	MonthDay struct {
		Month uint8 `json:"month"`
		Day   uint8 `json:"day"`
	}
)

func (h *Handler) loadTransactions(c *gin.Context) error {
	transactionDetails, err := h.TransactionsSrv.LoadTransactions(c.Request.Context(), 1)
	if err != nil {
		return err
	}

	dto := toDTO(transactionDetails)
	c.JSON(200, dto)
	return nil
}

func toDTO(transactions []transactions.TransactionDetail) []TransactionDetail {
	transactionsDTO := make([]TransactionDetail, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsDTO = append(transactionsDTO, TransactionDetail{
			ID: transaction.ID,
			Date: MonthDay{
				Month: transaction.Date.Month,
				Day:   transaction.Date.Day,
			},
			TransactionAmount: transaction.TransactionAmount,
		})
	}

	return transactionsDTO
}
