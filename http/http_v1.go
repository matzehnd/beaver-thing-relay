package http

import (
	"beaver/thing-relay/socket"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerService struct {
	socketService *socket.SocketService
}

func NewV1Handler(router *gin.RouterGroup, socketService socket.SocketService) {
	handler := &HandlerService{socketService: &socketService}
	router.POST("/things/:id", handler.thingsAction)
	router.GET("/things", handler.getThings)
}

type ActionTO struct {
	Action string  `json:"action"`
	Item   string  `json:"item"`
	Value  *string `json:"value,omitempty"`
}

func (h *HandlerService) thingsAction(c *gin.Context) {
	id := c.Param("id")
	var payload ActionTO
	if err := c.ShouldBindJSON((&payload)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.socketService.SendJson(id, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *HandlerService) getThings(c *gin.Context) {
	c.JSON(http.StatusOK, h.socketService.GetConnectionIds())
}
