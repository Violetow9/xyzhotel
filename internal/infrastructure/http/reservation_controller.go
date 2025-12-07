package http

import (
	"net/http"
	"time"
	"xyzhotel/internal/application"
	"xyzhotel/internal/domain/customer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type makeReservation struct {
	CustomerID     customer.ID `json:"customer_id"`
	CheckInDate    string      `json:"check_in_date"`
	AmountOfNights int         `json:"amount_of_nights"`
	Rooms          []string    `json:"rooms"`
}

type ReservationController struct {
	BookRoom           *application.BookRoomHandler
	ConfirmReservation *application.ConfirmReservationHandler
	CancelReservation  *application.CancelReservationHandler
	CheckoutRoom       *application.CheckoutRoomHandler
}

func (c *ReservationController) MakeReservationHandler(ctx *gin.Context) {
	var in makeReservation
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	checkInDate, err := time.Parse("2006-01-02", in.CheckInDate)

	cmd := &application.BookRoomCmd{
		CustomerID:     in.CustomerID,
		CheckInDate:    checkInDate,
		AmountOfNights: in.AmountOfNights,
		RoomNumbers:    in.Rooms,
	}

	res, err := c.BookRoom.Handle(ctx, cmd)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"id": res,
	})
}

func (c *ReservationController) ConfirmReservationHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	resID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	cmd := &application.ConfirmReservationCmd{ReservationID: resID}
	if err := c.ConfirmReservation.Handle(ctx, cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "confirmed"})
}

func (c *ReservationController) CancelReservationHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	resID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	cmd := &application.CancelReservationCmd{ReservationID: resID}
	if err := c.CancelReservation.Handle(ctx, cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (c *ReservationController) CheckoutRoomHandler(ctx *gin.Context) {
	idParam := ctx.Param("id")
	resID, err := uuid.Parse(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid reservation id"})
		return
	}

	cmd := &application.CheckoutRoomCmd{ReservationID: resID}
	if err := c.CheckoutRoom.Handle(ctx, cmd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "checked out"})
}
