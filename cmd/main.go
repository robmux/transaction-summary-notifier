package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/robmux/transaction-summary-notifier/pkg/configs"
	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
	"github.com/robmux/transaction-summary-notifier/pkg/repositories"
	"github.com/robmux/transaction-summary-notifier/pkg/services/rest"
)

func main() {
	r := gin.Default()

	transactionSrv := transactions.New()
	handler := rest.New(transactionSrv)

	rest.MountRoutes(r, handler)

	r.POST("/mails/notifications", func(c *gin.Context) {
		config := configs.GetMailConfig()
		em := repositories.NewSender(config)

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
