package application

import (
	"context"
	"errors"
	"xyzhotel/internal/domain/reservation"
)

type RoomHistoryQuery struct {
	RoomNumber string
}

type RoomHistoryHandler struct {
	ReservationRepo reservation.Repository
}

func (h RoomHistoryHandler) Handle(ctx context.Context, query RoomHistoryQuery) ([]reservation.Reservation, error) {
	if query.RoomNumber == "" {
		return nil, errors.New("room number is required")
	}

	return h.ReservationRepo.FindHistoryByRoom(ctx, query.RoomNumber)
}
