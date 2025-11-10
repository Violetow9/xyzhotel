package reservation

import (
	"time"
	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/room"

	"github.com/google/uuid"
)

type ID = uuid.UUID

type State = string

type CheckInDate = time.Time

type AmountOfNights = uint8

const (
	Pending   State = "PENDING"
	Confirmed State = "CONFIRMED"
	Cancelled State = "CANCELLED"
)

type Reservation struct {
	ID             ID
	Customer       *customer.Customer
	CheckInDate    CheckInDate
	AmountOfNights AmountOfNights
	Rooms          []*room.Room
	State          State
}
