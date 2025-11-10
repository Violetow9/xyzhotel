package room

import "context"

type Repository interface {
	Save(ctx context.Context, room Room) error

	Update(ctx context.Context, room Room) error

	Delete(ctx context.Context, id ID) error

	FindByID(ctx context.Context, id ID) (*Room, error)

	ListAll(ctx context.Context) []*Room
}
