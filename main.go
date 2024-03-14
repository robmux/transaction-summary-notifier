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
		c.String(200, "loading transactions")
	})

	r.POST("/notify", makeHTTPHandler(notifyUsers))

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
