package application

import (
	"context"
	"errors"
	"time"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/reservation"
	"xyzhotel/internal/domain/room"

	"github.com/google/uuid"
)

type MockCustomerRepo struct {
	Customers map[uuid.UUID]*customer.Customer
}

func NewMockCustomerRepo() *MockCustomerRepo {
	return &MockCustomerRepo{
		Customers: make(map[uuid.UUID]*customer.Customer),
	}
}

func (m *MockCustomerRepo) GetCustomerByID(ctx context.Context, id customer.ID) (*customer.Customer, error) {
	if c, ok := m.Customers[id]; ok {
		return c, nil
	}
	return nil, errors.New("customer not found")
}

func (m *MockCustomerRepo) GetCustomerByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	for _, c := range m.Customers {
		if c.Email == email {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockCustomerRepo) CreateCustomer(ctx context.Context, c *customer.Customer) error {
	if _, ok := m.Customers[c.ID]; ok {
		return errors.New("customer already exists")
	}
	m.Customers[c.ID] = c
	return nil
}

func (m *MockCustomerRepo) UpdateCustomer(ctx context.Context, c *customer.Customer) error {
	m.Customers[c.ID] = c
	return nil
}

func (m *MockCustomerRepo) ListCustomers(ctx context.Context) ([]customer.Customer, error) {
	var list []customer.Customer
	for _, c := range m.Customers {
		list = append(list, *c)
	}
	return list, nil
}

func (m *MockCustomerRepo) DeleteCustomer(ctx context.Context, id customer.ID) error {
	delete(m.Customers, id)
	return nil
}

type MockRoomRepo struct {
	Rooms map[string]*room.Room
}

func NewMockRoomRepo() *MockRoomRepo {
	return &MockRoomRepo{
		Rooms: make(map[string]*room.Room),
	}
}

func (m *MockRoomRepo) FindByID(ctx context.Context, id string) (*room.Room, error) {
	if r, ok := m.Rooms[id]; ok {
		return r, nil
	}
	return nil, errors.New("room not found")
}

func (m *MockRoomRepo) ListAll(ctx context.Context) ([]*room.Room, error) {
	var list []*room.Room
	for _, r := range m.Rooms {
		list = append(list, r)
	}
	return list, nil
}

func (m *MockRoomRepo) Save(ctx context.Context, r room.Room) error {
	return nil
}

func (m *MockRoomRepo) Update(ctx context.Context, r room.Room) error {
	return nil
}

func (m *MockRoomRepo) Delete(ctx context.Context, id room.ID) error {
	return nil
}

type MockReservationRepo struct {
	Reservations      map[uuid.UUID]*reservation.Reservation
	ForceAvailability bool
}

func NewMockReservationRepo() *MockReservationRepo {
	return &MockReservationRepo{
		Reservations:      make(map[uuid.UUID]*reservation.Reservation),
		ForceAvailability: true,
	}
}

func (m *MockReservationRepo) IsRoomAvailable(ctx context.Context, roomID string, date time.Time, nights int) (bool, error) {
	return m.ForceAvailability, nil
}

func (m *MockReservationRepo) Save(ctx context.Context, r *reservation.Reservation) error {
	m.Reservations[r.ID] = r
	return nil
}

func (m *MockReservationRepo) FindByID(ctx context.Context, id reservation.ID) (*reservation.Reservation, error) {
	if r, ok := m.Reservations[id]; ok {
		return r, nil
	}
	return nil, errors.New("reservation not found")
}

func (m *MockReservationRepo) Update(ctx context.Context, r *reservation.Reservation) error {
	m.Reservations[r.ID] = r
	return nil
}

func (m *MockReservationRepo) ListAll(ctx context.Context) ([]reservation.Reservation, error) {
	var list []reservation.Reservation
	for _, r := range m.Reservations {
		list = append(list, *r)
	}
	return list, nil
}

func (m *MockReservationRepo) FindHistoryByRoom(ctx context.Context, roomID string) ([]reservation.Reservation, error) {
	var list []reservation.Reservation
	for _, r := range m.Reservations {
		if r.RoomNumber == roomID {
			list = append(list, *r)
		}
	}
	return list, nil
}

func (m *MockReservationRepo) Delete(ctx context.Context, id reservation.ID) error {
	return nil
}
