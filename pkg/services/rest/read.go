package rest

import "github.com/gin-gonic/gin"

func (h *Handler) GetTransactionsSummary(c *gin.Context) error {
	summary, err := h.SummarySrv.GetSummary(c.Request.Context())
	if err != nil {
		return err
	}

	c.JSON(200, summary)
	return nil
}
