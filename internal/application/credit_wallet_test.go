package application

import (
	"context"
	"testing"
	customer2 "xyzhotel/internal/domain/customer"
	money2 "xyzhotel/internal/domain/money"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func createCreditWallet(initialMoney float64, moneyToCredit float64, currencyTarget money2.Currency) (*customer2.Customer, *CreditWalletCmd, *CreditWalletHandler) {
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
	cmd := &CreditWalletCmd{
		CustomerID: cust.ID,
		Money: money2.Money{
			Amount:   decimal.NewFromFloat(moneyToCredit),
			Currency: currencyTarget,
		},
	}
	handler := &CreditWalletHandler{
		CustomerService: customerService,
	}
	return cust, cmd, handler
}

func TestCreditWalletHandler_InSameCurrency(t *testing.T) {
	// Given a credit wallet handler
	// And a customer with a wallet
	// And an amount to be credited in the same currency as the wallet
	cust, cmd, handler := createCreditWallet(100, 10, money2.EUR)

	// When processing the credit
	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then verify the wallet balance is updated correctly
	expectedBalance := decimal.NewFromInt(110)
	if !cust.Wallet.Balance().Equal(expectedBalance) {
		t.Errorf("expected wallet balance to be %s, got %s", expectedBalance.String(), cust.Wallet.Balance().String())
	}
}

func TestCreditWalletHandler_InOtherCurrency(t *testing.T) {
	// Given a credit wallet handler
	// And a customer with a wallet
	// And an amount to be credited in a different currency than the wallet
	initialMoney := 100.0
	moneyToCredit := 10.0
	cust, cmd, handler := createCreditWallet(initialMoney, moneyToCredit, money2.USD)

	// When processing the credit
	err := handler.Handle(context.Background(), cmd)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then verify the wallet balance is updated correctly
	expectedBalance := initialMoney + (moneyToCredit * money2.USDToEURRate.InexactFloat64())
	if !cust.Wallet.Balance().Equal(decimal.NewFromFloat(expectedBalance)) {
		t.Errorf("expected wallet balance to be %.2f, got %s", expectedBalance, cust.Wallet.Balance().String())
	}
}
