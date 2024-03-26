package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/summary"
	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
	"github.com/robmux/transaction-summary-notifier/pkg/repositories"
	"github.com/robmux/transaction-summary-notifier/pkg/services/rest"
)

func main() {
	r := gin.Default()

	transactionsLoader := repositories.New()
	transactionSrv := transactions.New(transactionsLoader)

	summarySrv := summary.New(transactionSrv)

	handler := rest.New(transactionSrv, summarySrv)

	rest.MountRoutes(r, handler)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	err := r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Println("error ", err.Error())
	}
}
