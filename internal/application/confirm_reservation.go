package application

import (
	"context"
	"fmt"
	"xyzhotel/internal/domain/money"
	"xyzhotel/internal/domain/reservation"
)

type ConfirmReservationCmd struct {
	ReservationID reservation.ID
}

type ConfirmReservationHandler struct {
	ReservationRepo reservation.Repository
	DebitWallet     *DebitWalletHandler
}

func (h ConfirmReservationHandler) Handle(ctx context.Context, cmd *ConfirmReservationCmd) error {
	res, err := h.ReservationRepo.FindByID(ctx, cmd.ReservationID)
	if err != nil {
		return fmt.Errorf("reservation not found: %w", err)
	}

	remainingCents := res.DueAmount()
	if remainingCents <= 0 {
		return nil // Déjà payé
	}

	debitCmd := &DebitWalletCmd{
		CustomerID: res.CustomerID,
		Money: money.Money{
			AmountCents: remainingCents,
			Currency:    money.EUR,
		},
	}

	if err := h.DebitWallet.Handle(ctx, debitCmd); err != nil {
		return fmt.Errorf("payment failed: %w", err)
	}

	if err := res.Confirm(remainingCents); err != nil {
		return err
	}

	if err := h.ReservationRepo.Update(ctx, res); err != nil {
		return fmt.Errorf("failed to update reservation: %w", err)
	}

	return nil
}
