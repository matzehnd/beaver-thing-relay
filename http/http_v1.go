package http

import (
	"beaver/thing-relay/socket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerService struct {
	socketService *socket.SocketService
}

func NewV1Handler(router *gin.RouterGroup, socketService socket.SocketService) {
	handler := &HandlerService{socketService: &socketService}
	router.POST("/things/:id", handler.thingsAction)
}

type RegisterUserTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *HandlerService) thingsAction(c *gin.Context) {
	id := c.Param("id")
	log.Println(id)
	err := h.socketService.SendJson(id, gin.H{"hallo": "test"})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
