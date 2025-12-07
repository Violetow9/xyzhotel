package room

import "context"

type Repository interface {
	Save(ctx context.Context, room Room) error

	FindByID(ctx context.Context, id string) (*Room, error)

	ListAll(ctx context.Context) ([]*Room, error)
}
