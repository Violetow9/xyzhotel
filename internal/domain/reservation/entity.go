package reservation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidNights = errors.New("a reservation must be at least 1 night")
	ErrAlreadyPaid   = errors.New("reservation is already confirmed")
)

type ID = uuid.UUID
type CustomerID = uuid.UUID
type RoomNumber = string

type State string

const (
	Pending   State = "PENDING"
	Confirmed State = "CONFIRMED"
	Cancelled State = "CANCELLED"
	Completed State = "COMPLETED"
)

type Reservation struct {
	ID               ID
	CustomerID       CustomerID
	RoomNumber       RoomNumber
	CheckInDate      time.Time
	AmountOfNights   int
	State            State
	TotalAmountCents int
	PaidAmountCents  int
}

func NewReservation(
	customerID CustomerID,
	roomNumber RoomNumber,
	roomPriceCents int,
	checkIn time.Time,
	nights int,
) (*Reservation, error) {
	if nights < 1 {
		return nil, ErrInvalidNights
	}

	totalAmount := roomPriceCents * nights
	deposit := totalAmount / 2

	return &Reservation{
		ID:               uuid.New(),
		CustomerID:       customerID,
		RoomNumber:       roomNumber,
		CheckInDate:      checkIn,
		AmountOfNights:   nights,
		State:            Pending,
		TotalAmountCents: totalAmount,
		PaidAmountCents:  deposit,
	}, nil
}

func (r *Reservation) Confirm(paymentAmountCents int) error {
	if r.State == Confirmed {
		return ErrAlreadyPaid
	}

	remaining := r.TotalAmountCents - r.PaidAmountCents
	if paymentAmountCents < remaining {
		return errors.New("insufficient amount to confirm reservation")
	}

	r.PaidAmountCents += paymentAmountCents
	r.State = Confirmed
	return nil
}

func (r *Reservation) IsActive() bool {
	return r.State == Pending || r.State == Confirmed
}

func (r *Reservation) DueAmount() int {
	return r.TotalAmountCents - r.PaidAmountCents
}
