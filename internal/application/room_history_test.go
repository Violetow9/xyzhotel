package application

import (
	"context"
	"testing"
	"time"

	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

func TestRoomHistoryHandler_EmptyRoomNumber(t *testing.T) {
	handler := &RoomHistoryHandler{}
	query := RoomHistoryQuery{RoomNumber: ""}

	_, err := handler.Handle(context.Background(), query)

	if err == nil {
		t.Error("expected error for empty room number")
	}
}

func TestRoomHistoryHandler_Success(t *testing.T) {
	repo := NewMockReservationRepo()
	ctx := context.Background()

	// Reservations for room 101
	res1 := &reservation.Reservation{
		ID:          uuid.New(),
		RoomNumber:  "101",
		CheckInDate: time.Now().AddDate(0, 0, -5),
		State:       reservation.Completed,
	}
	res2 := &reservation.Reservation{
		ID:          uuid.New(),
		RoomNumber:  "101",
		CheckInDate: time.Now().AddDate(0, 0, 5),
		State:       reservation.Confirmed,
	}
	repo.Save(ctx, res1)
	repo.Save(ctx, res2)

	// Reservation for room 102 (should be ignored)
	res3 := &reservation.Reservation{
		ID:         uuid.New(),
		RoomNumber: "102",
	}
	repo.Save(ctx, res3)

	handler := &RoomHistoryHandler{ReservationRepo: repo}
	query := RoomHistoryQuery{RoomNumber: "101"}

	history, err := handler.Handle(ctx, query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("expected 2 reservations in history, got %d", len(history))
	}

	for _, res := range history {
		if res.RoomNumber != "101" {
			t.Errorf("expected room 101, got %s", res.RoomNumber)
		}
	}
}
