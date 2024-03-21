package rest

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/robmux/transaction-summary-notifier/pkg/domains/transactions"
)

type (
	TransactionsManager interface {
		LoadTransactions(ctx context.Context, userID uint64) ([]transactions.TransactionDetail, error)
	}
)

type Handler struct {
	TransactionsSrv TransactionsManager
}

func New(txManager TransactionsManager) *Handler {
	return &Handler{
		TransactionsSrv: txManager,
	}
}

func MountRoutes(r *gin.Engine, handler *Handler) {
	// Read
	r.GET("/ping", makeHTTPHandler(handler.Ping))

	// Write

	r.POST("/load-transactions", makeHTTPHandler(handler.loadTransactions))

}

func (h *Handler) Ping(ctx *gin.Context) error {
	ctx.String(200, "pong")
	return nil
}

type appHandler func(ctx *gin.Context) error

// makeHTTPHandlers converts an handler with err to a gin handler
func makeHTTPHandler(h appHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h(ctx)
		if err != nil {
			ctx.String(500, err.Error())
		}

		// Nothing here, because h should have already sent the response
	}
}
