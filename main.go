package main

import (
	"beaver/thing-relay/http"
	"beaver/thing-relay/socket"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("no port defined")
	}
	IdpUrl := os.Getenv("IDP_URL")

	if IdpUrl == "" {
		panic("no IdpUrl defined")
	}

	socketService := socket.NewSocketService()

	r := gin.Default()
	r.Use(TokenCheck(IdpUrl))
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/ws", socket.ConnectionHandler(socketService))
	v1 := r.Group("/v1")
	http.NewV1Handler(v1, *socketService)

	r.Run(":" + port)
}
