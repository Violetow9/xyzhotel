package application

import (
	"context"
	"errors"
	"testing"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"

	"github.com/google/uuid"
)

func TestDebitWalletHandler_CustomerNotFound(t *testing.T) {
	repo := NewMockCustomerRepo()
	converter := money.NewConverter()

	handler := &DebitWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	cmd := &DebitWalletCmd{
		CustomerID: uuid.New(),
		Money:      money.New(100, money.EUR),
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when customer not found")
	}
}

func TestDebitWalletHandler_ConversionError(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", customer.ZeroWallet())
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &DebitWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	cmd := &DebitWalletCmd{
		CustomerID: id,
		Money: money.Money{
			AmountCents: 1000,
			Currency:    "UNKNOWN",
		},
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error for invalid currency")
	}
}

func TestDebitWalletHandler_InsufficientFunds(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()

	// Wallet 10.00 EUR
	wallet := customer.NewWallet(1000)
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", wallet)
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &DebitWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	// Try to debit 20.00 EUR
	cmd := &DebitWalletCmd{
		CustomerID: id,
		Money:      money.New(2000, money.EUR),
	}

	err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}

func TestDebitWalletHandler_Success(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()

	// Wallet 100.00 EUR
	wallet := customer.NewWallet(10000)
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", wallet)
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &DebitWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	// Debit 40.00 EUR
	cmd := &DebitWalletCmd{
		CustomerID: id,
		Money:      money.New(4000, money.EUR),
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedCust, _ := repo.GetCustomerByID(context.Background(), id)

	// 100.00 - 40.00 = 60.00 EUR (6000 cents)
	if updatedCust.Wallet.Balance() != 6000 {
		t.Errorf("expected balance 6000, got %d", updatedCust.Wallet.Balance())
	}
}

func TestDebitWalletHandler_SuccessWithConversion(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()

	// Wallet 100.00 EUR
	wallet := customer.NewWallet(10000)
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", wallet)
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &DebitWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	// Debit 100.00 USD -> 85.00 EUR (Rate 0.85)
	cmd := &DebitWalletCmd{
		CustomerID: id,
		Money:      money.New(10000, money.USD),
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedCust, _ := repo.GetCustomerByID(context.Background(), id)

	// 10000 - 8500 = 1500
	if updatedCust.Wallet.Balance() != 1500 {
		t.Errorf("expected balance 1500, got %d", updatedCust.Wallet.Balance())
	}
}
