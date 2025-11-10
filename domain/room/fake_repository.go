package room

import "context"

var _ Repository = &FakeRepository{}

type FakeRepository struct {
	Rooms map[ID]Room
}

func (f *FakeRepository) Save(ctx context.Context, room Room) error {
	f.Rooms[room.ID] = room
	return nil
}

func (f *FakeRepository) Update(ctx context.Context, room Room) error {
	f.Rooms[room.ID] = room
	return nil
}

func (f *FakeRepository) Delete(ctx context.Context, id ID) error {
	delete(f.Rooms, id)
	return nil
}

func (f *FakeRepository) FindByID(ctx context.Context, id ID) (*Room, error) {
	room, exists := f.Rooms[id]
	if !exists {
		return nil, nil
	}
	return &room, nil
}

func (f *FakeRepository) ListAll(ctx context.Context) []*Room {
	rooms := make([]*Room, 0, len(f.Rooms))
	for _, room := range f.Rooms {
		r := room
		rooms = append(rooms, &r)
	}
	return rooms
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		Rooms: make(map[ID]Room),
	}
}
