package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello, Gin!")
	})

	r.POST("/load-transactions", func(c *gin.Context) {
		csvFile, err := loadFile("input/user_1_transactions.csv")
		if err != nil {
			c.JSON(500, err)
			return
		}

		columns := []string{
			"TxID", "Date", "TransactionAmount",
		}
		transactions, err := readCSV(columns, csvFile)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		c.String(200, fmt.Sprintf("%+v", transactions))
	})

	r.POST("/notify", makeHTTPHandler(notifyUsers))

	r.GET("/transactions/summary", func(c *gin.Context) {
		csvFile, err := loadFile("input/user_1_transactions.csv")
		if err != nil {
			c.JSON(500, err)
			return
		}

		columns := []string{
			"TxID", "Date", "TransactionAmount",
		}
		transactions, err := readCSV(columns, csvFile)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		avg := GetAverageCreditAndDebit(transactions)

		c.JSON(200, fmt.Sprintf("%+v", avg))
	})

	r.GET("/transactions/summary/avg", func(c *gin.Context) {
		csvFile, err := loadFile("input/user_1_transactions.csv")
		if err != nil {
			c.JSON(500, err)
			return
		}

		columns := []string{
			"TxID", "Date", "TransactionAmount",
		}
		transactions, err := readCSV(columns, csvFile)
		if err != nil {
			c.JSON(500, err.Error())
			return
		}

		avg := GetAverageDebit(transactions)
		avgCredit := GetAverageCredit(transactions)

		c.JSON(200, fmt.Sprintf("Average debit \n %+v \n Average credit \n %+v", avg, avgCredit))
	})

	err := r.Run(":3000")
	if err != nil {
		fmt.Println("error ", err.Error())
	}
}

type handler func(ctx *gin.Context) error

func makeHTTPHandler(h handler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h(ctx)
		if err != nil {
			ctx.String(500, err.Error())
		}

		// Nothing here, because h should have already sent the response
	}
}

func notifyUsers(ctx *gin.Context) error {
	ctx.String(200, "users notified")
	return nil
}
