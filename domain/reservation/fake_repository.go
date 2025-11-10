package reservation

import "context"

var _ Repository = &FakeRepository{}

type FakeRepository struct {
	reservations map[ID]*Reservation
}

func (f *FakeRepository) Save(ctx context.Context, reservation *Reservation) error {
	f.reservations[reservation.ID] = reservation
	return nil
}

func (f *FakeRepository) Update(ctx context.Context, reservation *Reservation) error {
	f.reservations[reservation.ID] = reservation
	return nil
}

func (f *FakeRepository) Delete(ctx context.Context, id ID) error {
	delete(f.reservations, id)
	return nil
}

func (f *FakeRepository) FindByID(ctx context.Context, id ID) (*Reservation, error) {
	reservation, exists := f.reservations[id]
	if !exists {
		return nil, nil
	}
	return reservation, nil
}

func (f *FakeRepository) ListAll(ctx context.Context) []Reservation {
	reservations := make([]Reservation, 0, len(f.reservations))
	for _, reservation := range f.reservations {
		reservations = append(reservations, *reservation)
	}
	return reservations
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		reservations: make(map[ID]*Reservation),
	}
}
