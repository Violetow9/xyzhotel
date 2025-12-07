package application

import (
	"context"
	"testing"
	"time"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"
	"xyzhotel/internal/domain/reservation"

	"github.com/google/uuid"
)

func TestConfirmReservationHandler_ReservationNotFound(t *testing.T) {
	resRepo := NewMockReservationRepo()
	custRepo := NewMockCustomerRepo()

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &ConfirmReservationHandler{
		ReservationRepo: resRepo,
		DebitWallet:     debitHandler,
	}

	cmd := &ConfirmReservationCmd{
		ReservationID: uuid.New(),
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when reservation not found")
	}
}

func TestConfirmReservationHandler_AlreadyPaid(t *testing.T) {
	resID := uuid.New()
	custID := uuid.New()

	// 100 EUR total, 100 EUR paid
	res := &reservation.Reservation{
		ID:               resID,
		CustomerID:       custID,
		TotalAmountCents: 10000,
		PaidAmountCents:  10000,
		State:            reservation.Confirmed,
	}

	resRepo := NewMockReservationRepo()
	resRepo.Save(context.Background(), res)

	// Wallet 0
	wallet := customer.NewWallet(0)
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", wallet)
	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &ConfirmReservationHandler{
		ReservationRepo: resRepo,
		DebitWallet:     debitHandler,
	}

	cmd := &ConfirmReservationCmd{
		ReservationID: resID,
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Balance should remain 0 (no debit attempted)
	updatedCust, _ := custRepo.GetCustomerByID(context.Background(), custID)
	if updatedCust.Wallet.Balance() != 0 {
		t.Error("wallet should not have been debited")
	}
}

func TestConfirmReservationHandler_InsufficientFunds(t *testing.T) {
	resID := uuid.New()
	custID := uuid.New()

	// 100 EUR total, 50 EUR paid. 50 EUR remaining.
	res := &reservation.Reservation{
		ID:               resID,
		CustomerID:       custID,
		TotalAmountCents: 10000,
		PaidAmountCents:  5000,
		State:            reservation.Pending,
	}

	resRepo := NewMockReservationRepo()
	resRepo.Save(context.Background(), res)

	// Wallet only 10 EUR
	wallet := customer.NewWallet(1000)
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", wallet)
	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &ConfirmReservationHandler{
		ReservationRepo: resRepo,
		DebitWallet:     debitHandler,
	}

	cmd := &ConfirmReservationCmd{
		ReservationID: resID,
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error due to insufficient funds")
	}
}

func TestConfirmReservationHandler_Success(t *testing.T) {
	resID := uuid.New()
	custID := uuid.New()

	// 100 EUR total, 50 EUR paid. 50 EUR remaining.
	res := &reservation.Reservation{
		ID:               resID,
		CustomerID:       custID,
		TotalAmountCents: 10000,
		PaidAmountCents:  5000,
		State:            reservation.Pending,
		CheckInDate:      time.Now(),
		AmountOfNights:   1,
	}

	resRepo := NewMockReservationRepo()
	resRepo.Save(context.Background(), res)

	// Wallet 100 EUR
	wallet := customer.NewWallet(10000)
	cust := customer.NewCustomer(custID, "John", "john@test.com", "123", wallet)
	custRepo := NewMockCustomerRepo()
	custRepo.CreateCustomer(context.Background(), cust)

	debitHandler := &DebitWalletHandler{
		CustomerRepo: custRepo,
		Converter:    money.NewConverter(),
	}

	handler := &ConfirmReservationHandler{
		ReservationRepo: resRepo,
		DebitWallet:     debitHandler,
	}

	cmd := &ConfirmReservationCmd{
		ReservationID: resID,
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedRes, _ := resRepo.FindByID(context.Background(), resID)
	if updatedRes.State != reservation.Confirmed {
		t.Errorf("expected state Confirmed, got %s", updatedRes.State)
	}
	if updatedRes.PaidAmountCents != 10000 {
		t.Errorf("expected paid amount 10000, got %d", updatedRes.PaidAmountCents)
	}

	// Check Wallet (100 - 50 = 50 EUR)
	updatedCust, _ := custRepo.GetCustomerByID(context.Background(), custID)
	if updatedCust.Wallet.Balance() != 5000 {
		t.Errorf("expected balance 5000, got %d", updatedCust.Wallet.Balance())
	}
}
