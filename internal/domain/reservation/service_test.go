package reservation

import (
	"testing"
	_ "testing"
	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/room"
)

var _ Repository = &InMemoryRepository{}

type InMemoryRepository struct {
	reservations map[ID]*Reservation
}

func (r InMemoryRepository) Save(reservation *Reservation) error {
	r.reservations[reservation.ID] = reservation
	return nil
}

func (r InMemoryRepository) Update(reservation *Reservation) error {
	r.reservations[reservation.ID] = reservation
	return nil
}

func (r InMemoryRepository) Delete(id ID) error {
	delete(r.reservations, id)
	return nil
}

func (r InMemoryRepository) FindByID(id ID) (*Reservation, error) {
	return r.reservations[id], nil
}

func (r InMemoryRepository) ListAll() ([]Reservation, error) {
	reservations := make([]Reservation, 0, len(r.reservations))
	for _, reservation := range r.reservations {
		reservations = append(reservations, *reservation)
	}
	return reservations, nil
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		reservations: make(map[ID]*Reservation),
	}
}

// Show reservation history from a room ID
func TestShowReservationHistory(t *testing.T) {
	// Given a reservation service with a mock ReservationRepository
	// And a room ID
	// And some reservations associated with that room ID in the ReservationRepository
	service := Service{
		Repository: NewInMemoryRepository(),
	}
	rooms := []room.Room{
		room.New(101, room.Standard),
		room.New(102, room.Deluxe),
	}
	customers := []customer.Customer{
		customer.NewCustomer(1, "Alice", "alice@gmail.com", "1234567890"),
		customer.NewCustomer(2, "Bob", "bob@gmail.com", "0987654321"),
	}
	reservations := []Reservation{
		New(1, customers[0], "2025-11-06", 2, []room.Room{rooms[0]}),
		New(2, customers[1], "2025-11-07", 3, []room.Room{rooms[1]}),
		New(3, customers[0], "2025-11-08", 1, []room.Room{rooms[0], rooms[1]}),
	}
	for i := range reservations {
		err := service.Repository.Save(&reservations[i])
		if err != nil {
			t.Fatalf("Failed to save reservation: %v", err)
		}
	}

	// When ShowReservationHistory is called with the room ID
	history := service.ShowReservationHistory(101)

	// Then it should return the correct reservation history for that room ID
	if len(history) != 2 {
		t.Errorf("Expected 2 reservations in history, got %d", len(history))
	}
}
