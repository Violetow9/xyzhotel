package reservation

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, reservation *Reservation) error

	Update(ctx context.Context, reservation *Reservation) error

	FindByID(ctx context.Context, id ID) (*Reservation, error)

	ListAll(ctx context.Context) ([]Reservation, error)

	IsRoomAvailable(ctx context.Context, roomID string, date time.Time, nights int) (bool, error)

	FindHistoryByRoom(ctx context.Context, roomID string) ([]Reservation, error)
}
