package application

import (
	"context"
	"testing"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"

	"github.com/google/uuid"
)

func TestCreditWalletHandler_CustomerNotFound(t *testing.T) {
	repo := NewMockCustomerRepo()
	converter := money.NewConverter()

	handler := &CreditWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	cmd := &CreditWalletCmd{
		CustomerID: uuid.New(),
		Money:      money.New(100, money.EUR),
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error when customer not found")
	}
}

func TestCreditWalletHandler_ConversionError(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", customer.ZeroWallet())
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &CreditWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	cmd := &CreditWalletCmd{
		CustomerID: id,
		Money: money.Money{
			AmountCents: 1000,
			Currency:    "INVALID",
		},
	}

	err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error for invalid currency")
	}
}

func TestCreditWalletHandler_SuccessEUR(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()

	initialWallet := customer.NewWallet(1000) // 10.00 EUR
	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", initialWallet)
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &CreditWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	cmd := &CreditWalletCmd{
		CustomerID: id,
		Money:      money.New(5000, money.EUR), // 50.00 EUR
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedCust, _ := repo.GetCustomerByID(context.Background(), id)
	expectedBalance := 6000 // 1000 + 5000
	if updatedCust.Wallet.Balance() != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, updatedCust.Wallet.Balance())
	}
}

func TestCreditWalletHandler_SuccessUSD(t *testing.T) {
	repo := NewMockCustomerRepo()
	id := uuid.New()

	cust := customer.NewCustomer(id, "Test", "test@test.com", "123", customer.ZeroWallet())
	repo.CreateCustomer(context.Background(), cust)

	converter := money.NewConverter()

	handler := &CreditWalletHandler{
		CustomerRepo: repo,
		Converter:    converter,
	}

	// 100 USD -> 85 EUR (Rate 0.85)
	cmd := &CreditWalletCmd{
		CustomerID: id,
		Money:      money.New(10000, money.USD), // 100.00 USD
	}

	err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedCust, _ := repo.GetCustomerByID(context.Background(), id)

	// 10000 * 0.85 = 8500
	expectedBalance := 8500
	if updatedCust.Wallet.Balance() != expectedBalance {
		t.Errorf("expected balance %d, got %d", expectedBalance, updatedCust.Wallet.Balance())
	}
}
