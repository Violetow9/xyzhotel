package http

import (
	"net/http"
	"xyzhotel/internal/application"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	ListRooms *application.ListRoomsHandler
}

func (c *RoomController) ListRoomsHandler(ctx *gin.Context) {
	rooms, err := c.ListRooms.Handle(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, rooms)
}
