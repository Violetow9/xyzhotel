package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"
	"xyzhotel/internal/domain/reservation"
	"xyzhotel/internal/domain/room"
)

type BookRoomCmd struct {
	CustomerID     customer.ID
	CheckInDate    time.Time
	AmountOfNights int
	RoomNumbers    []string
}

var (
	ErrCheckInDateInPast  = errors.New("check-in date cannot be in the past")
	ErrNoRoomsSelected    = errors.New("no rooms selected for booking")
	ErrAmountOfNightsZero = errors.New("amount of nights must be greater than zero")
)

type BookRoomHandler struct {
	CustomerRepository customer.Repository
	ReservationRepo    reservation.Repository
	RoomRepository     room.Repository
	DebitWallet        *DebitWalletHandler
}

func (h BookRoomHandler) Handle(ctx context.Context, cmd *BookRoomCmd) ([]reservation.ID, error) {
	if len(cmd.RoomNumbers) == 0 {
		return nil, ErrNoRoomsSelected
	}
	if cmd.AmountOfNights <= 0 {
		return nil, ErrAmountOfNightsZero
	}

	// Normaliser les dates à minuit pour éviter les problèmes de comparaison
	today := truncateToDay(time.Now())
	checkIn := truncateToDay(cmd.CheckInDate)

	if checkIn.Before(today) {
		return nil, ErrCheckInDateInPast
	}

	if _, err := h.CustomerRepository.GetCustomerByID(ctx, cmd.CustomerID); err != nil {
		return nil, fmt.Errorf("customer invalid: %w", err)
	}

	var roomsToBook []room.Room
	totalDepositCents := 0

	for _, roomNum := range cmd.RoomNumbers {
		r, err := h.RoomRepository.FindByID(ctx, roomNum)
		if err != nil {
			return nil, fmt.Errorf("room %s not found: %w", roomNum, err)
		}

		isAvailable, err := h.ReservationRepo.IsRoomAvailable(ctx, string(r.ID), checkIn, cmd.AmountOfNights)
		if err != nil {
			return nil, fmt.Errorf("availability check failed: %w", err)
		}
		if !isAvailable {
			return nil, fmt.Errorf("room %s is not available for these dates", r.ID)
		}

		roomsToBook = append(roomsToBook, *r)

		config := r.GetConfig()
		totalPriceRoom := config.PriceCents * cmd.AmountOfNights

		// Acompte de 50%
		depositRoom := totalPriceRoom / 2
		totalDepositCents += depositRoom
	}

	debitCmd := &DebitWalletCmd{
		CustomerID: cmd.CustomerID,
		Money: money.Money{
			AmountCents: totalDepositCents,
			Currency:    money.EUR,
		},
	}

	if err := h.DebitWallet.Handle(ctx, debitCmd); err != nil {
		return nil, err
	}

	var createdReservationIDs []reservation.ID

	for _, r := range roomsToBook {
		res, err := reservation.NewReservation(
			cmd.CustomerID,
			string(r.ID),
			r.GetConfig().PriceCents,
			checkIn,
			cmd.AmountOfNights,
		)
		if err != nil {
			return nil, err
		}

		if err := h.ReservationRepo.Save(ctx, res); err != nil {
			return nil, fmt.Errorf("failed to save reservation: %w", err)
		}
		createdReservationIDs = append(createdReservationIDs, res.ID)
	}

	return createdReservationIDs, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
