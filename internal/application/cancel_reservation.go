package application

import (
	"context"
	"fmt"
	"xyzhotel/internal/domain/reservation"
)

type CancelReservationCmd struct {
	ReservationID reservation.ID
}

type CancelReservationHandler struct {
	ReservationRepo reservation.Repository
}

func (h CancelReservationHandler) Handle(ctx context.Context, cmd *CancelReservationCmd) error {
	res, err := h.ReservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	res.State = reservation.Cancelled

	if err := h.ReservationRepo.Update(ctx, res); err != nil {
		return fmt.Errorf("failed to cancel reservation: %w", err)
	}

	return nil
}
