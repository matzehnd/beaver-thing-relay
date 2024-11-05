package socket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConnectionHandler(s *SocketService) func(*gin.Context) {
	return func(gin *gin.Context) {
		conn, err := upgrader.Upgrade(gin.Writer, gin.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}
		defer conn.Close()
		s.RegisterConnection("1", conn)

	}
}
