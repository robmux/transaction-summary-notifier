package summary

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
)

type (
	Srv struct {
		transactionsLoader TransactionsLoader
	}

	TransactionsLoader interface {
		LoadTransactions(ctx context.Context, userID uint64) ([]transactions.TransactionDetail, error)
	}
)

const (
	DebitBalance  = "Debit"
	CreditBalance = "Credit"
)

func New(transactionsLoader TransactionsLoader) *Srv {
	return &Srv{
		transactionsLoader: transactionsLoader,
	}
}

func (s *Srv) GetSummary(ctx context.Context, userID uint64) (*GeneralSummary, error) {
	details, err := s.transactionsLoader.LoadTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalBalance := s.GetTotalBalanceInAccount(ctx, details)

	balanceType := DebitBalance
	if totalBalance.IsPositive() {
		balanceType = CreditBalance
	}
	if totalBalance.IsZero() {
		balanceType = ""
	}

	numTxsByMonthMap := s.GetNumberOfTransactionsGroupedByMonth(ctx, details)
	numTxByMonth := make([]TransactionsByMonth, 0, len(numTxsByMonthMap))
	for _, month := range numTxsByMonthMap {
		numTxByMonth = append(numTxByMonth, month)
	}

	averagesByMonth := s.GetAverageCreditAndDebit(ctx, details)

	avgDebit := s.GetAverageDebit(ctx, details)
	avgCredit := s.GetAverageCredit(ctx, details)

	summary := GeneralSummary{
		TotalBalance: AmountDetail{
			Amount:     totalBalance,
			AmountType: balanceType,
		},
		NumberTransactionsByMonth: numTxByMonth,
		AveragesByMonth:           averagesByMonth,
		AverageCredit:             avgDebit,
		AverageDebit:              avgCredit,
	}
	return &summary, nil
}

func (s *Srv) GetTotalBalanceInAccount(
	ctx context.Context,
	transactions []transactions.TransactionDetail) decimal.Decimal {
	totalBalance := decimal.NewFromFloat(0.0)

	for i := 0; i < len(transactions); i++ {
		totalBalance = totalBalance.Add(transactions[i].TransactionAmount)
	}

	return totalBalance
}

func (s *Srv) GetNumberOfTransactionsGroupedByMonth(
	ctx context.Context,
	transactions []transactions.TransactionDetail) map[uint8]TransactionsByMonth {
	byMonthMap := make(map[uint8]TransactionsByMonth, 12)
	for _, transaction := range transactions {
		// Get the value and if it's not present then the default is returned
		monthData := byMonthMap[transaction.Date.Month]

		// Make sure to add the month in case it was the first time and the default value was returned
		monthData.Month = transaction.Date.Month
		monthData.TransactionsQuantity += uint64(1)

		byMonthMap[transaction.Date.Month] = monthData
	}

	return byMonthMap
}

// GetAverageCreditAndDebit does not include zeros in the avg
func (s *Srv) GetAverageCreditAndDebit(
	ctx context.Context,
	transactions []transactions.TransactionDetail) AveragesByMonth {
	// Create map to keep the average
	byMonthMap := make(map[uint8]AverageByMonth, 12)

	avgByMonth := make(map[string]avgData)
	for idx := range transactions {
		avgType := DebitBalance

		avgKey := fmt.Sprintf("%d", transactions[idx].Date.Month)

		avg := avgByMonth[avgKey]
		avg.count++
		avg.totalSum = avg.totalSum.Add(transactions[idx].TransactionAmount)

		// Update average
		avg.avg = avg.totalSum.Div(decimal.NewFromInt(avg.count))
		if avg.avg.IsPositive() {
			avgType = CreditBalance
		}
		if avg.avg.IsZero() {
			avgType = ""
		}

		// update the map
		avgByMonth[avgKey] = avg

		byMonth := byMonthMap[transactions[idx].Date.Month]
		byMonth.Month = transactions[idx].Date.Month
		byMonth.Average = avg.avg
		byMonth.TransactionType = avgType
		byMonthMap[transactions[idx].Date.Month] = byMonth
	}

	monthsData := make([]AverageByMonth, 0, len(byMonthMap))
	for _, month := range byMonthMap {
		monthsData = append(monthsData, month)
	}

	return AveragesByMonth{
		AvgsByMonth: monthsData,
	}
}

func (s *Srv) GetAverageDebit(
	ctx context.Context,
	transactions []transactions.TransactionDetail) AmountDetail {
	debitCounter := 0
	total := decimal.NewFromFloat(0.0)
	for _, transaction := range transactions {
		if transaction.TransactionAmount.IsNegative() {
			debitCounter++
			total = total.Add(transaction.TransactionAmount)
		}
	}

	detail := AmountDetail{
		AmountType: DebitBalance,
		Amount:     decimal.NewFromFloat(0.0),
	}

	if debitCounter == 0 {
		return detail
	}

	detail.Amount = total.Div(decimal.NewFromInt(int64(debitCounter)))
	return detail
}

func (s *Srv) GetAverageCredit(
	ctx context.Context,
	transactions []transactions.TransactionDetail) AmountDetail {
	debitCounter := 0
	total := decimal.NewFromFloat(0.0)
	for _, transaction := range transactions {
		if transaction.TransactionAmount.IsPositive() {
			debitCounter++
			total = total.Add(transaction.TransactionAmount)
		}
	}

	detail := AmountDetail{
		AmountType: "Credit",
		Amount:     decimal.NewFromFloat(0.0),
	}

	if debitCounter == 0 {
		return detail
	}

	detail.Amount = total.Div(decimal.NewFromInt(int64(debitCounter)))
	return detail
}

func makeSliceFromMap(mapFrom map[uint8]AverageByMonth, to []AverageByMonth) []AverageByMonth {
	for _, average := range mapFrom {
		to = append(to, average)
	}

	return to
}
