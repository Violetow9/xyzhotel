package application

import (
	"context"
	"errors"
	"testing"
	"time"
	"xyzhotel/domain/customer"
	"xyzhotel/domain/reservation"
	"xyzhotel/domain/room"

	"github.com/google/uuid"
)

func TestBookRoomHandler_CheckInDateInPast(t *testing.T) {
	// Given a book room handler
	// And a check-in date in the past
	ctx := context.Background()
	repository := customer.NewFakeRepository()
	bookRoomHandler := &BookRoomHandler{
		CustomerRepository: repository,
		ReservationService: nil,
		RoomRepository:     nil,
	}

	// When booking a room
	_, err := bookRoomHandler.Handle(ctx, &BookRoomCmd{
		CustomerID:     customer.ID{},
		CheckInDate:    time.Now().AddDate(0, 0, -1), // Yesterday
		AmountOfNights: 2,
		Rooms:          []room.ID{},
	})

	// Then expect an ErrCheckInDateInPast error
	if !errors.Is(err, ErrCheckInDateInPast) {
		t.Errorf("expected ErrCheckInDateInPast, but got: %v", err)
	}
}

func TestBookRoomHandler_NoRoomsSelected(t *testing.T) {
	// Given a book room handler
	// And no rooms selected
	ctx := context.Background()
	repository := customer.NewFakeRepository()
	bookRoomHandler := &BookRoomHandler{
		CustomerRepository: repository,
		ReservationService: nil,
		RoomRepository:     nil,
	}

	// When booking a room
	_, err := bookRoomHandler.Handle(ctx, &BookRoomCmd{
		CustomerID:     customer.ID{},
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 2,
		Rooms:          []room.ID{},
	})

	// Then expect an ErrNoRoomsSelected error
	if !errors.Is(err, ErrNoRoomsSelected) {
		t.Errorf("expected ErrNoRoomsSelected, but got: %v", err)
	}
}

func TestBookRoomHandler_AmountOfNightsZero(t *testing.T) {
	// Given a book room handler
	// And an amount of nights equal to zero
	ctx := context.Background()
	repository := customer.NewFakeRepository()
	bookRoomHandler := &BookRoomHandler{
		CustomerRepository: repository,
		ReservationService: nil,
		RoomRepository:     nil,
	}

	// When booking a room
	_, err := bookRoomHandler.Handle(ctx, &BookRoomCmd{
		CustomerID:     customer.ID{},
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 0,
		Rooms:          []room.ID{uuid.New()},
	})

	// Then expect an ErrAmountOfNightsZero error
	if !errors.Is(err, ErrAmountOfNightsZero) {
		t.Errorf("expected ErrAmountOfNightsZero, but got: %v", err)
	}
}

func TestBookRoomHandler_NonExistingCustomer(t *testing.T) {
	// Given a book room handler
	// And a customer ID that does not exist
	ctx := context.Background()
	repository := customer.NewFakeRepository()
	bookRoomHandler := &BookRoomHandler{
		CustomerRepository: repository,
		ReservationService: nil,
		RoomRepository:     nil,
	}

	// When booking a room
	_, err := bookRoomHandler.Handle(ctx, &BookRoomCmd{
		CustomerID:     uuid.New(),
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 2,
		Rooms:          []room.ID{uuid.New()},
	})

	// Then expect an ErrCustomerNotFound error
	if !errors.Is(err, customer.ErrCustomerNotFound) {
		t.Errorf("expected ErrCustomerNotFound, but got: %v", err)
	}
}

func TestBookRoomHandler_OccupiedRooms(t *testing.T) {
	// Given a book room handler
	// And rooms that are already booked
	ctx := context.Background()
	customerRepository := customer.NewFakeRepository()
	customerA := &customer.Customer{
		ID: uuid.New(),
	}
	customerB := &customer.Customer{
		ID: uuid.New(),
	}
	err := customerRepository.CreateCustomer(ctx, customerA)
	if err != nil {
		t.Fatalf("failed to set up customer A: %v", err)
	}
	err = customerRepository.CreateCustomer(ctx, customerB)
	if err != nil {
		t.Fatalf("failed to set up customer B: %v", err)
	}

	roomRepository := room.NewFakeRepository()
	reservationRepository := reservation.NewFakeRepository()
	reservationService := &reservation.Service{
		ReservationRepository: reservationRepository,
		RoomRepository:        roomRepository,
	}
	CheckInDate := time.Now().AddDate(0, 0, 1)
	roomID := uuid.New()
	err = reservationService.CreateReservation(ctx, &reservation.Reservation{
		ID:             uuid.New(),
		Customer:       customerA,
		CheckInDate:    CheckInDate,
		AmountOfNights: 2,
		Rooms: []*room.Room{
			{ID: roomID},
		},
	})
	if err != nil {
		t.Fatalf("failed to set up existing reservation: %v", err)
	}

	bookRoomHandler := &BookRoomHandler{
		CustomerRepository: customerRepository,
		ReservationService: reservationService,
		RoomRepository:     roomRepository,
	}

	// When booking a room
	var res *reservation.Reservation
	res, err = bookRoomHandler.Handle(ctx, &BookRoomCmd{
		CustomerID:     customerB.ID,
		CheckInDate:    CheckInDate,
		AmountOfNights: 2,
		Rooms:          []room.ID{roomID},
	})

	// Then expect an ErrRoomAlreadyBooked error
	if !errors.Is(err, ErrRoomAlreadyBooked) {
		t.Errorf("expected ErrRoomAlreadyBooked, but got: %v", err)
	}

	if res == nil {
		t.Errorf("expected no reservation to be created, but got one")
	}
}
