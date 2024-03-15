package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// loadFile loads the file completely by putting it into memory.
func loadFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type TransactionDetail struct {
	ID   uint64
	Date MonthDay

	TransactionAmount decimal.Decimal
}

type MonthDay struct {
	Month uint8
	Day   uint8
}

func readCSV(columns []string, csvBytes []byte) ([]TransactionDetail, error) {
	csvReader := csv.NewReader(bytes.NewReader(csvBytes))
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	transactions := make([]TransactionDetail, 0, len(data))
	for i, line := range data {
		if i == 0 {
			// validate headers are the expected
			continue
		}

		if len(line) != len(columns) {
			return nil, fmt.Errorf("line %d has %d columns but should have exactly %d", i, len(line), len(columns))
		}

		transactionX := TransactionDetail{}
		for j, column := range line {
			switch columns[j] {
			case "TxID":
				transactionX.ID, err = strconv.ParseUint(column, 10, 64)
				if err != nil {
					return nil, err
				}
			case "Date":
				dayMonthStrs := strings.Split(column, "/")
				if len(dayMonthStrs) != 2 {
					return nil, fmt.Errorf("invalid date values %s", column)
				}
				month, err := strconv.ParseUint(dayMonthStrs[0], 10, 64)
				if err != nil {
					return nil, err
				}
				day, err := strconv.ParseUint(dayMonthStrs[1], 10, 64)
				if err != nil {
					return nil, err
				}

				if day > 31 {
					return nil, fmt.Errorf("invalid day number %d", day)
				}
				if month > 12 {
					return nil, fmt.Errorf("invalid month number %d", month)
				}

				transactionX.Date = MonthDay{
					Day:   uint8(day),
					Month: uint8(month),
				}

			case "TransactionAmount":
				if len(column) <= 1 {
					return nil, fmt.Errorf("invalid transaction amount %s", column)
				}

				if column[0] != '-' && column[0] != '+' {
					return nil, fmt.Errorf("invalid transaction amount, does not have sign %s", column)
				}

				amountFloat, err := strconv.ParseFloat(column, 64)
				if err != nil {
					return nil, err
				}
				transactionX.TransactionAmount = decimal.NewFromFloat(amountFloat)
			}
		}

		transactions = append(transactions, transactionX)
	}

	return transactions, nil
}
