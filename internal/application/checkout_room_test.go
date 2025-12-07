package application

import (
	"context"
	"testing"

	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

func TestCheckoutRoomHandler_ReservationNotFound(t *testing.T) {
	repo := NewMockReservationRepo()
	handler := &CheckoutRoomHandler{ReservationRepo: repo}

	cmd := &CheckoutRoomCmd{
		ReservationID: uuid.New(),
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when reservation not found")
	}
}

func TestCheckoutRoomHandler_Success(t *testing.T) {
	repo := NewMockReservationRepo()
	resID := uuid.New()

	existingRes := &reservation.Reservation{
		ID:    resID,
		State: reservation.Confirmed,
	}
	repo.Save(context.Background(), existingRes)

	handler := &CheckoutRoomHandler{ReservationRepo: repo}

	cmd := &CheckoutRoomCmd{
		ReservationID: resID,
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedRes, _ := repo.FindByID(context.Background(), resID)
	if updatedRes.State != reservation.Completed {
		t.Errorf("expected state %s, got %s", reservation.Completed, updatedRes.State)
	}
}
