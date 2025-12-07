package application

import (
	"context"
	"testing"
	"time"

	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

func midnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func TestOccupiedRoomsHandler_NoReservations(t *testing.T) {
	repo := NewMockReservationRepo()
	handler := &OccupiedRoomsHandler{ReservationRepo: repo}

	query := OccupiedRoomsQuery{}
	list, err := handler.Handle(context.Background(), query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected 0 occupied rooms, got %d", len(list))
	}
}

func TestOccupiedRoomsHandler_NoActiveReservations(t *testing.T) {
	repo := NewMockReservationRepo()
	ctx := context.Background()
	now := midnight(time.Now())

	pastRes := &reservation.Reservation{
		ID:             uuid.New(),
		RoomNumber:     "101",
		CheckInDate:    now.AddDate(0, 0, -10),
		AmountOfNights: 2,
		State:          reservation.Confirmed,
	}
	repo.Save(ctx, pastRes)

	futureRes := &reservation.Reservation{
		ID:             uuid.New(),
		RoomNumber:     "102",
		CheckInDate:    now.AddDate(0, 0, 10),
		AmountOfNights: 2,
		State:          reservation.Confirmed,
	}
	repo.Save(ctx, futureRes)

	cancelledRes := &reservation.Reservation{
		ID:             uuid.New(),
		RoomNumber:     "103",
		CheckInDate:    now.AddDate(0, 0, -1),
		AmountOfNights: 5,
		State:          reservation.Cancelled,
	}
	repo.Save(ctx, cancelledRes)

	handler := &OccupiedRoomsHandler{ReservationRepo: repo}
	query := OccupiedRoomsQuery{}

	list, err := handler.Handle(ctx, query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected 0 occupied rooms, got %d", len(list))
	}
}

func TestOccupiedRoomsHandler_ActiveReservations(t *testing.T) {
	repo := NewMockReservationRepo()
	ctx := context.Background()
	now := midnight(time.Now())

	res1 := &reservation.Reservation{
		ID:             uuid.New(),
		RoomNumber:     "101",
		CheckInDate:    now,
		AmountOfNights: 2,
		State:          reservation.Pending,
	}
	repo.Save(ctx, res1)

	res2 := &reservation.Reservation{
		ID:             uuid.New(),
		RoomNumber:     "102",
		CheckInDate:    now.AddDate(0, 0, -1),
		AmountOfNights: 3,
		State:          reservation.Confirmed,
	}
	repo.Save(ctx, res2)

	handler := &OccupiedRoomsHandler{ReservationRepo: repo}
	query := OccupiedRoomsQuery{}

	list, err := handler.Handle(ctx, query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 occupied rooms, got %d", len(list))
	}
}
