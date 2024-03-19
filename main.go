package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
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

	r.POST("/mails/notifications", func(c *gin.Context) {
		config := getMailConfig()
		em := NewSender(config)

		err := em.SendEmailNotification()
		if err != nil {
			c.JSON(500, err.Error())
			return
		}
	})

	err := r.Run(":3000")
	if err != nil {
		fmt.Println("error ", err.Error())
	}
}

func getMailConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("MAIL_USER")
	if len(user) == 0 {
		panic("mail user missing")
	}

	password := os.Getenv("MAIL_PASSWORD")
	host := os.Getenv("MAIL_HOST")
	serverAddr := os.Getenv("MAIL_SERVER_ADDR")

	return Config{
		userMail: user,
		password: password,

		host:       host,
		serverAddr: serverAddr,
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
