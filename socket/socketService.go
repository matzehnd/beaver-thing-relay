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
	var isClosed = false
	connection.SetCloseHandler(func(code int, text string) error {
		s.DeRegister(thingId)
		isClosed = true
		return nil
	})
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			if isClosed {
				break
			}
			log.Println("Read error:", err)
			break
		}
		fmt.Printf("Received: %s\n", message)
	}
}

func (s *SocketService) DeRegister(thingId string) {
	delete(s.connections, thingId)
}

func (s *SocketService) GetConnectionIds() []string {
	keys := make([]string, 0, len(s.connections))
	for k := range s.connections {
		keys = append(keys, k)
	}
	return keys
}

func (s *SocketService) SendJson(thingId string, message interface{}) error {
	connection, exists := s.connections[thingId]
	if !exists {
		return fmt.Errorf("unable to find connection to thing: %s", thingId)
	}

	err := connection.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("unable to send Message: %s", err)
	}

	return nil
}
