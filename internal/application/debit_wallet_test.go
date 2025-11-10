package application

import (
	"context"
	"errors"
	"testing"
	customer2 "xyzhotel/internal/domain/customer"
	money2 "xyzhotel/internal/domain/money"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func createDebitWallet(initialMoney float64, moneyToDebit float64, currencyTarget money2.Currency) (*customer2.Customer, *DebitWalletCmd, *DebitWallerHandler) {
	cust := customer2.NewCustomer(
		uuid.New(),
		"John",
		"john@gmail.com",
		"0987654321",
		customer2.NewWallet(decimal.NewFromFloat(initialMoney)),
	)
	customerRepository := customer2.NewFakeRepositoryWithCustomers(cust)
	customerService := &customer2.Service{
		Repository: customerRepository,
	}
	cmd := &DebitWalletCmd{
		CustomerID: cust.ID,
		Money: money2.Money{
			Amount:   decimal.NewFromFloat(moneyToDebit),
			Currency: currencyTarget,
		},
	}
	handler := &DebitWallerHandler{
		CustomerService: customerService,
	}
	return cust, cmd, handler
}

func TestDebitWalletHandler_InSameCurrency(t *testing.T) {
	// Given a debit wallet handler
	// And a customer with a wallet
	// And an amount to be debited in the same currency as the wallet
	cust, cmd, handler := createDebitWallet(100, 10, money2.EUR)

	// When processing the credit
	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then verify the wallet balance is updated correctly
	expectedBalance := decimal.NewFromInt(90)
	if !cust.Wallet.Balance().Equal(expectedBalance) {
		t.Errorf("expected wallet balance to be %s, got %s", expectedBalance.String(), cust.Wallet.Balance().String())
	}
}

func TestDebitWalletHandler_InOtherCurrency(t *testing.T) {
	// Given a debit wallet handler
	// And a customer with a wallet
	// And an amount to be debited in a different currency than the wallet
	initialMoney := 100.0
	amountToDebit := 50.0
	cust, cmd, handler := createDebitWallet(initialMoney, amountToDebit, money2.USD)

	// When processing the debit
	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then verify the wallet balance is updated correctly
	expectedBalance := decimal.NewFromFloat(initialMoney - (amountToDebit * money2.USDToEURRate.InexactFloat64()))
	if !cust.Wallet.Balance().Equal(expectedBalance) {
		t.Errorf("expected wallet balance to be %s, got %s", expectedBalance, cust.Wallet.Balance().String())
	}
}

func TestDebitWalletHandler_InsufficientFunds(t *testing.T) {
	// Given a debit wallet handler
	// And a customer with a wallet
	// And an amount to be paid in the same currency as the wallet
	// But the wallet has insufficient funds
	initialMoney := 10.0
	moneyToDebit := 50.0
	cust, cmd, handler := createDebitWallet(initialMoney, moneyToDebit, money2.EUR)

	// When processing the debit
	err := handler.Handle(context.Background(), cmd)
	if err == nil {
		t.Errorf("expected error due to insufficient funds, got nil")
	}

	// Then it should return an error indicating insufficient funds
	// and the wallet balance remains unchanged
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected insufficient funds error, got: %v", err)
	}

	expectedBalance := decimal.NewFromFloat(moneyToDebit)
	if cust.Wallet.Balance().Equal(expectedBalance) {
		t.Errorf("expected wallet balance to be %s, got %s", expectedBalance.String(), cust.Wallet.Balance().String())
	}
}
