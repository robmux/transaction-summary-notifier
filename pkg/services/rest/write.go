package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
)

type (
	TransactionDetailDTO struct {
		ID   uint64      `json:"id"`
		Date MonthDayDTO `json:"date"`

		TransactionAmount decimal.Decimal `json:"transaction_amount"`
	}

	MonthDayDTO struct {
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

func toDTO(transactions []transactions.TransactionDetail) []TransactionDetailDTO {
	transactionsDTO := make([]TransactionDetailDTO, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsDTO = append(transactionsDTO, TransactionDetailDTO{
			ID: transaction.ID,
			Date: MonthDayDTO{
				Month: transaction.Date.Month,
				Day:   transaction.Date.Day,
			},
			TransactionAmount: transaction.TransactionAmount,
		})
	}

	return transactionsDTO
}
