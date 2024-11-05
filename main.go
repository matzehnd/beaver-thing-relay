package main

import (
	"beaver/thing-relay/config"
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

	dbConn := os.Getenv("DB")

	if dbConn == "" {
		panic("no db connection string")
	}

	cfg := config.LoadConfig(dbConn)
	config.InitDb(*cfg)
	defer cfg.DB.Close()

	socketService := socket.NewSocketService()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	r.GET("/ws", socket.ConnectionHandler(socketService))
	v1 := r.Group("/v1", TokenCheck(IdpUrl))
	http.NewV1Handler(v1, *socketService)

	r.Run(":" + port)
}
