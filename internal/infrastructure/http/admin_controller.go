package http

import (
	"net/http"
	"xyzhotel/internal/application"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	RoomHistory   *application.RoomHistoryHandler
	OccupiedRooms *application.OccupiedRoomsHandler
}

func (c *AdminController) HistoryHandler(ctx *gin.Context) {
	roomID := ctx.Param("roomNumber")
	query := application.RoomHistoryQuery{
		RoomNumber: roomID,
	}

	history, err := c.RoomHistory.Handle(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, history)
}

func (c *AdminController) OccupiedHandler(ctx *gin.Context) {
	query := application.OccupiedRoomsQuery{}

	occupied, err := c.OccupiedRooms.Handle(ctx, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, occupied)
}
