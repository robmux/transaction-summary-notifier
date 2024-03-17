package main

import (
	"fmt"
	"github.com/shopspring/decimal"
)

type TransactionsByMonth struct {
	Month                uint8
	TransactionsQuantity uint64
}

func GetTotalBalanceInAccount(transactions []TransactionDetail) decimal.Decimal {
	totalBalance := decimal.NewFromFloat(0.0)

	for i := 0; i < len(transactions); i++ {
		totalBalance.Add(transactions[i].TransactionAmount)
	}

	return totalBalance
}

func GetNumberOfTransactionsGroupedByMonth(transactions []TransactionDetail) map[uint8]TransactionsByMonth {
	byMonthMap := make(map[uint8]TransactionsByMonth, 12)
	for _, transaction := range transactions {
		monthData := byMonthMap[transaction.Date.Month]
		monthData.TransactionsQuantity += uint64(1)

		byMonthMap[transaction.Date.Month] = monthData
	}

	return byMonthMap
}

type AverageByMonth struct {
	Month   uint8
	Average decimal.Decimal

	TransactionType string
}

type avgData struct {
	count    int64
	totalSum decimal.Decimal

	avg decimal.Decimal
}

type AveragesByMonth struct {
	Debit  []AverageByMonth
	Credit []AverageByMonth
}

// GetAverageCreditAndDebit does not include zeros in the avg
func GetAverageCreditAndDebit(transactions []TransactionDetail) AveragesByMonth {
	byMonthMapDebit := make(map[uint8]AverageByMonth, 12)
	byMonthMapCredit := make(map[uint8]AverageByMonth, 12)

	avgByMonthAndType := make(map[string]avgData)
	credit := "Credit"
	debit := "Debit"
	for idx := range transactions {
		avgType := debit
		if transactions[idx].TransactionAmount.IsPositive() {
			avgType = credit
		}

		avgKey := fmt.Sprintf("%d_type_%s", transactions[idx].Date.Month, avgType)

		avg := avgByMonthAndType[avgKey]
		avg.count++
		avg.totalSum = avg.totalSum.Add(transactions[idx].TransactionAmount)

		// Update average
		avg.avg = avg.totalSum.Div(decimal.NewFromInt(avg.count))
		// update the map
		avgByMonthAndType[avgKey] = avg

		// update result
		if avgType == credit {
			byMonth := byMonthMapCredit[transactions[idx].Date.Month]
			byMonth.Average = avg.avg
			byMonth.TransactionType = credit
			byMonthMapCredit[transactions[idx].Date.Month] = byMonth
			continue
		}

		byMonth := byMonthMapDebit[transactions[idx].Date.Month]
		byMonth.Average = avg.avg
		byMonth.TransactionType = debit
		byMonthMapDebit[transactions[idx].Date.Month] = byMonth

	}

	debitMonths := make([]AverageByMonth, 0, len(byMonthMapDebit))
	creditMonths := make([]AverageByMonth, 0, len(byMonthMapCredit))

	debitMonths = makeSliceFromMap(byMonthMapDebit, debitMonths)
	creditMonths = makeSliceFromMap(byMonthMapCredit, creditMonths)

	return AveragesByMonth{
		Debit:  debitMonths,
		Credit: creditMonths,
	}
}

type AmountDetail struct {
	Amount     decimal.Decimal
	AmountType string
}

func GetAverageDebit(transactions []TransactionDetail) AmountDetail {
	debitCounter := 0
	total := decimal.NewFromFloat(0.0)
	for _, transaction := range transactions {
		if transaction.TransactionAmount.IsNegative() {
			debitCounter++
			total = total.Add(transaction.TransactionAmount)
		}
	}

	detail := AmountDetail{
		AmountType: "Debit",
		Amount:     decimal.NewFromFloat(0.0),
	}

	if debitCounter == 0 {
		return detail
	}

	detail.Amount = total.Div(decimal.NewFromInt(int64(debitCounter)))
	return detail
}

func GetAverageCredit(transactions []TransactionDetail) AmountDetail {
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
