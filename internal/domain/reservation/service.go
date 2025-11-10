package reservation

import (
	"context"
	"time"
	"xyzhotel/internal/domain/room"
)

type Service struct {
	ReservationRepository Repository
	RoomRepository        room.Repository
}

func (rs *Service) ShowReservationHistory(ctx context.Context, roomId room.ID) []Reservation {
	reservations := rs.ReservationRepository.ListAll(ctx)
	var history []Reservation
	for _, reservation := range reservations {
		for _, r := range reservation.Rooms {
			if r.ID == roomId {
				history = append(history, reservation)
				break
			}
		}
	}
	return history
}

func (rs *Service) CreateReservation(ctx context.Context, reservation *Reservation) error {
	return rs.ReservationRepository.Save(ctx, reservation)
}

func (rs *Service) IsRoomAvailable(ctx context.Context, roomID room.ID, startDate time.Time, amountOfNights uint8) (bool, error) {
	reservations := rs.ReservationRepository.ListAll(ctx)
	endDate := startDate.AddDate(0, 0, int(amountOfNights))
	for _, reservation := range reservations {
		for _, r := range reservation.Rooms {
			if r.ID == roomID {
				resStart := reservation.CheckInDate
				resEnd := resStart.AddDate(0, 0, int(reservation.AmountOfNights))
				if startDate.Before(resEnd) && endDate.After(resStart) {
					return false, nil
				}
			}
		}
	}
	return true, nil
}

func (rs *Service) ListAvailableRooms(ctx context.Context) []*room.Room {
	rooms := rs.RoomRepository.ListAll(ctx)
	var availableRooms []*room.Room
	for _, r := range rooms {
		isAvailable, _ := rs.IsRoomAvailable(ctx, r.ID, time.Now(), 1)
		if isAvailable {
			availableRooms = append(availableRooms, r)
		}
	}
	return availableRooms
}

func (rs *Service) ListOccupiedRooms(ctx context.Context) []*room.Room {
	rooms := rs.RoomRepository.ListAll(ctx)
	var occupiedRooms []*room.Room
	for _, r := range rooms {
		isAvailable, _ := rs.IsRoomAvailable(ctx, r.ID, time.Now(), 1)
		if !isAvailable {
			occupiedRooms = append(occupiedRooms, r)
		}
	}
	return occupiedRooms
}
