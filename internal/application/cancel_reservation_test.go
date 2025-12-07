package application

import (
	"context"
	"testing"
	"time"

	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

func TestCancelReservationHandler_ReservationNotFound(t *testing.T) {
	repo := NewMockReservationRepo()
	handler := &CancelReservationHandler{ReservationRepo: repo}

	cmd := &CancelReservationCmd{
		ReservationID: uuid.New(),
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when reservation not found")
	}
}

func TestCancelReservationHandler_Success(t *testing.T) {
	repo := NewMockReservationRepo()
	resID := uuid.New()

	existingRes := &reservation.Reservation{
		ID:             resID,
		State:          reservation.Pending,
		CheckInDate:    time.Now(),
		AmountOfNights: 2,
	}
	repo.Save(context.Background(), existingRes)

	handler := &CancelReservationHandler{ReservationRepo: repo}

	cmd := &CancelReservationCmd{
		ReservationID: resID,
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedRes, _ := repo.FindByID(context.Background(), resID)
	if updatedRes.State != reservation.Cancelled {
		t.Errorf("expected state %s, got %s", reservation.Cancelled, updatedRes.State)
	}
}
