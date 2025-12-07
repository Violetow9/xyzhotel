package application

import (
	"context"
	"errors"
	"testing"
	"time"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"
	"xyzhotel/internal/domain/room"

	"github.com/google/uuid"
)

func TestBookRoomHandler_CheckInDateInPast(t *testing.T) {
	handler := &BookRoomHandler{}
	cmd := &BookRoomCmd{
		CheckInDate:    time.Now().AddDate(0, 0, -1),
		AmountOfNights: 1,
		RoomNumbers:    []string{"101"},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrCheckInDateInPast) {
		t.Errorf("expected ErrCheckInDateInPast, got %v", err)
	}
}

func TestBookRoomHandler_NoRoomsSelected(t *testing.T) {
	handler := &BookRoomHandler{}
	cmd := &BookRoomCmd{
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 1,
		RoomNumbers:    []string{},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrNoRoomsSelected) {
		t.Errorf("expected ErrNoRoomsSelected, got %v", err)
	}
}

func TestBookRoomHandler_AmountOfNightsZero(t *testing.T) {
	handler := &BookRoomHandler{}
	cmd := &BookRoomCmd{
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 0,
		RoomNumbers:    []string{"101"},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrAmountOfNightsZero) {
		t.Errorf("expected ErrAmountOfNightsZero, got %v", err)
	}
}

func TestBookRoomHandler_CustomerNotFound(t *testing.T) {
	custRepo := NewMockCustomerRepo()
	handler := &BookRoomHandler{CustomerRepository: custRepo}

	cmd := &BookRoomCmd{
		CustomerID:     uuid.New(),
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 1,
		RoomNumbers:    []string{"101"},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error for non-existing customer")
	}
}

func TestBookRoomHandler_RoomNotAvailable(t *testing.T) {
	custID := uuid.New()
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", customer.ZeroWallet())

	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	roomRepo := NewMockRoomRepo()
	roomRepo.Rooms["101"] = &room.Room{ID: "101", Type: room.Standard}

	resRepo := NewMockReservationRepo()
	resRepo.ForceAvailability = false // Force l'indisponibilité

	handler := &BookRoomHandler{
		CustomerRepository: custRepo,
		RoomRepository:     roomRepo,
		ReservationRepo:    resRepo,
	}

	cmd := &BookRoomCmd{
		CustomerID:     custID,
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 1,
		RoomNumbers:    []string{"101"},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when room is not available")
	}
}

func TestBookRoomHandler_InsufficientFunds(t *testing.T) {
	custID := uuid.New()
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", customer.ZeroWallet())

	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	roomRepo := NewMockRoomRepo()
	roomRepo.Rooms["101"] = &room.Room{ID: "101", Type: room.Standard} // 50 EUR

	resRepo := NewMockReservationRepo()

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &BookRoomHandler{
		CustomerRepository: custRepo,
		RoomRepository:     roomRepo,
		ReservationRepo:    resRepo,
		DebitWallet:        debitHandler,
	}

	// Prix total 50 EUR, Acompte 25 EUR, Wallet 0 EUR
	cmd := &BookRoomCmd{
		CustomerID:     custID,
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 1,
		RoomNumbers:    []string{"101"},
	}

	_, err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestBookRoomHandler_Success(t *testing.T) {
	custID := uuid.New()
	wallet := customer.NewWallet(10000) // 100 EUR
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", wallet)

	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	roomRepo := NewMockRoomRepo()
	roomRepo.Rooms["101"] = &room.Room{ID: "101", Type: room.Standard} // 50 EUR/nuit

	resRepo := NewMockReservationRepo()

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &BookRoomHandler{
		CustomerRepository: custRepo,
		RoomRepository:     roomRepo,
		ReservationRepo:    resRepo,
		DebitWallet:        debitHandler,
	}

	// 1 nuit = 50 EUR, Acompte = 25 EUR.
	cmd := &BookRoomCmd{
		CustomerID:     custID,
		CheckInDate:    time.Now().AddDate(0, 0, 1),
		AmountOfNights: 1,
		RoomNumbers:    []string{"101"},
	}

	resIDs, err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resIDs) != 1 {
		t.Errorf("expected 1 reservation, got %d", len(resIDs))
	}

	if len(resRepo.Reservations) != 1 {
		t.Error("reservation was not saved in repo")
	}

	updatedCust, _ := custRepo.GetCustomerByID(context.Background(), custID)
	// 100 EUR - 25 EUR = 75 EUR (7500 cents)
	if updatedCust.Wallet.Balance() != 7500 {
		t.Errorf("expected wallet balance 7500, got %d", updatedCust.Wallet.Balance())
	}
}
