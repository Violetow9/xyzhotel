package application

import (
	"context"
	"fmt"
	"xyzhotel/internal/domain/reservation"
)

type CheckoutRoomCmd struct {
	ReservationID reservation.ID
}

type CheckoutRoomHandler struct {
	ReservationRepo reservation.Repository
}

func (h CheckoutRoomHandler) Handle(ctx context.Context, cmd *CheckoutRoomCmd) error {
	res, err := h.ReservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	res.State = reservation.Completed

	if err := h.ReservationRepo.Update(ctx, res); err != nil {
		return fmt.Errorf("failed to checkout reservation: %w", err)
	}

	return nil
}
