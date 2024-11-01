package main

import (
	"beaver/thing-relay/config"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("no port defined")
	}

	dbConn := os.Getenv("DB")

	if dbConn == "" {
		panic("no db connection string")
	}

	cfg := config.LoadConfig(dbConn)
	config.InitDb(*cfg)
	defer cfg.DB.Close()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	// v1 := r.Group("/v1")
	r.Run(":" + port)

}
