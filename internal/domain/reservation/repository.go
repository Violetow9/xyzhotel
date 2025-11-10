package reservation

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, reservation *Reservation) error

	Update(ctx context.Context, reservation *Reservation) error

	Delete(ctx context.Context, id ID) error

	FindByID(ctx context.Context, id ID) (*Reservation, error)

	ListAll(ctx context.Context) []Reservation
}
