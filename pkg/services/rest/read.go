package rest

import (
	"fmt"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/summary"
)

func (h *Handler) GetTransactionsSummary(c *gin.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user_id param %w", err)
	}

	summary, err := h.SummarySrv.GetSummary(c.Request.Context(), userID)
	if err != nil {
		return err
	}

	dtoSummary := toDTOSummary(summary)

	c.JSON(200, dtoSummary)
	return nil
}

type AveragesByMonthResponse struct {
	MonthsData []AverageByMonthResponse `json:"months_data"`
}

type TransactionsByMonthResponse struct {
	Month                uint8  `json:"month"`
	TransactionsQuantity uint64 `json:"transactions_quantity"`
}

type AverageByMonthResponse struct {
	Month   uint8           `json:"month"`
	Average decimal.Decimal `json:"average"`

	TransactionType string `json:"transaction_type"`
}

type GeneralSummaryResponse struct {
	TotalBalance AmountDetailResponse `json:"total_balance"`

	NumberTransactionsByMonth []TransactionsByMonthResponse `json:"number_transactions_by_month"`
	AveragesByMonth           AveragesByMonthResponse       `json:"averages_by_month"`

	AverageCredit AmountDetailResponse `json:"average_credit"`
	AverageDebit  AmountDetailResponse `json:"average_debit"`
}

type AmountDetailResponse struct {
	Amount     decimal.Decimal `json:"amount"`
	AmountType string          `json:"amount_type"`
}

func toDTOSummary(summary *summary.GeneralSummary) GeneralSummaryResponse {
	response := GeneralSummaryResponse{}

	totalBalance := AmountDetailResponse{
		Amount:     summary.TotalBalance.Amount,
		AmountType: summary.TotalBalance.AmountType,
	}
	response.TotalBalance = totalBalance

	txByMonth := make([]TransactionsByMonthResponse, 0, len(summary.NumberTransactionsByMonth))
	for _, txBy := range summary.NumberTransactionsByMonth {
		txByMonthDetail := TransactionsByMonthResponse{
			Month:                txBy.Month,
			TransactionsQuantity: txBy.TransactionsQuantity,
		}
		txByMonth = append(txByMonth, txByMonthDetail)

	}
	response.NumberTransactionsByMonth = txByMonth

	avgsByMonth := AveragesByMonthResponse{}
	monthsData := make([]AverageByMonthResponse, 0, len(summary.AveragesByMonth.AvgsByMonth))
	for _, byMonth := range summary.AveragesByMonth.AvgsByMonth {
		monthsData = append(monthsData, AverageByMonthResponse{
			Month:           byMonth.Month,
			Average:         byMonth.Average,
			TransactionType: byMonth.TransactionType,
		})
	}
	avgsByMonth.MonthsData = monthsData
	response.AveragesByMonth = avgsByMonth

	avgCredit := AmountDetailResponse{
		Amount:     summary.AverageCredit.Amount,
		AmountType: summary.AverageCredit.AmountType,
	}
	response.AverageCredit = avgCredit
	avgDebit := AmountDetailResponse{
		Amount:     summary.AverageDebit.Amount,
		AmountType: summary.AverageDebit.AmountType,
	}
	response.AverageDebit = avgDebit

	return response
}
