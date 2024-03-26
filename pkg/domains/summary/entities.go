package summary

import "github.com/shopspring/decimal"

type AmountDetail struct {
	Amount     decimal.Decimal
	AmountType string
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
	AvgsByMonth []AverageByMonth
}

type TransactionsByMonth struct {
	Month                uint8
	TransactionsQuantity uint64
}

type GeneralSummary struct {
	TotalBalance AmountDetail

	NumberTransactionsByMonth []TransactionsByMonth
	AveragesByMonth           AveragesByMonth

	AverageCredit AmountDetail
	AverageDebit  AmountDetail
}
