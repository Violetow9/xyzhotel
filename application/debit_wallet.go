package application

import (
	"context"
	"errors"

	"xyzhotel/domain/customer"
	"xyzhotel/domain/money"
)

type DebitWalletCmd struct {
	CustomerID customer.ID
	Money      money.Money
}

var (
	ErrInsufficientFunds = errors.New("not sufficient funds")
)

type DebitWallerHandler struct {
	CustomerService *customer.Service
}

func (h DebitWallerHandler) Handle(ctx context.Context, cmd *DebitWalletCmd) error {
	cust, err := h.CustomerService.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return err
	}

	converter, err := money.CurrencyToConverter(cmd.Money.Currency)
	if err != nil {
		return err
	}

	convertedAmount := converter.Convert(cmd.Money)
	if !cust.Wallet.HasSufficientFunds(convertedAmount) {
		return ErrInsufficientFunds
	}
	cust.Wallet.Debit(convertedAmount)
	return h.CustomerService.UpdateCustomer(ctx, cust)
}
