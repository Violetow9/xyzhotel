package application

import (
	"context"
	"time"
	"xyzhotel/internal/domain/reservation"
)

type OccupiedRoomsQuery struct{}

type OccupiedRoomsHandler struct {
	ReservationRepo reservation.Repository
}

func (h OccupiedRoomsHandler) Handle(ctx context.Context, query OccupiedRoomsQuery) ([]reservation.Reservation, error) {
	allReservations, err := h.ReservationRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var occupied []reservation.Reservation

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for _, res := range allReservations {
		isActiveStatus := res.State == reservation.Confirmed || res.State == reservation.Pending

		if isActiveStatus {
			checkOutDate := res.CheckInDate.AddDate(0, 0, res.AmountOfNights)

			isStayOngoing := (res.CheckInDate.Before(today) || res.CheckInDate.Equal(today)) && today.Before(checkOutDate)

			if isStayOngoing {
				occupied = append(occupied, res)
			}
		}
	}

	return occupied, nil
}
