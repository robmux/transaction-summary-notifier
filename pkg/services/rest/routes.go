package rest

import "github.com/gin-gonic/gin"

func MountRoutes(r *gin.Engine, handler *Handler) {
	// Read
	r.GET("/ping", makeHTTPHandler(handler.Ping))
	r.GET("/transactions/summary", makeHTTPHandler(handler.GetTransactionsSummary))

	// Write

	r.POST("/load-transactions", makeHTTPHandler(handler.loadTransactions))

}
