package rest

import "github.com/gin-gonic/gin"

func MountRoutes(r *gin.Engine, handler *Handler) {
	// Read
	r.GET("/ping", makeHTTPHandler(handler.Ping))
	r.GET("/users/:user_id/transactions/summary", makeHTTPHandler(handler.GetTransactionsSummary))

	// Write

	r.POST("/load-transactions", makeHTTPHandler(handler.loadTransactions))
	r.POST("/mails/notifications", makeHTTPHandler(handler.SendMails))
}
