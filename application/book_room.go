package application

import (
	"context"
	"errors"
	"time"
	"xyzhotel/domain/customer"
	"xyzhotel/domain/reservation"
	"xyzhotel/domain/room"

	"github.com/google/uuid"
)

type BookRoomCmd struct {
	CustomerID     customer.ID
	CheckInDate    time.Time
	AmountOfNights uint8
	Rooms          []room.ID
}

var (
	ErrRoomAlreadyBooked  = errors.New("room already booked for the selected dates")
	ErrCheckInDateInPast  = errors.New("check-in date cannot be in the past")
	ErrNoRoomsSelected    = errors.New("no rooms selected for booking")
	ErrAmountOfNightsZero = errors.New("amount of nights must be greater than zero")
)

type BookRoomHandler struct {
	CustomerRepository customer.Repository
	ReservationService *reservation.Service
	RoomRepository     room.Repository
}

func (h BookRoomHandler) Handle(ctx context.Context, cmd *BookRoomCmd) (*reservation.Reservation, error) {
	if cmd.CheckInDate.Before(time.Now()) {
		return nil, ErrCheckInDateInPast
	}

	if len(cmd.Rooms) == 0 {
		return nil, ErrNoRoomsSelected
	}

	if cmd.AmountOfNights == 0 {
		return nil, ErrAmountOfNightsZero
	}

	cust, err := h.CustomerRepository.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return nil, err
	}

	rooms := make([]*room.Room, 0, len(cmd.Rooms))
	for _, roomID := range cmd.Rooms {
		r, err := h.RoomRepository.FindByID(ctx, roomID)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}

	for _, r := range rooms {
		isAvailable, err := h.ReservationService.IsRoomAvailable(ctx, r.ID, cmd.CheckInDate, cmd.AmountOfNights)
		if err != nil {
			return nil, err
		}
		if !isAvailable {
			return nil, ErrRoomAlreadyBooked
		}
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	res := &reservation.Reservation{
		ID:             id,
		Customer:       cust,
		CheckInDate:    cmd.CheckInDate,
		AmountOfNights: cmd.AmountOfNights,
		Rooms:          rooms,
	}

	err = h.ReservationService.CreateReservation(ctx, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
