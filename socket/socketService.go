package socket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type SocketService struct {
	connections map[string]*websocket.Conn
}

func NewSocketService() *SocketService {
	return &SocketService{
		connections: make(map[string]*websocket.Conn),
	}
}

func (s *SocketService) RegisterConnection(thingId string, connection *websocket.Conn) {
	s.connections[thingId] = connection
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
	}
}

func (s *SocketService) SendJson(thingId string, message interface{}) error {
	connection, exists := s.connections[thingId]
	if !exists {
		return fmt.Errorf("unable to find connection to thing: %s", thingId)
	}

	err := connection.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("unable to send Message: %T", err)
	}

	return nil
}
