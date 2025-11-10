package application

import (
	"context"

	"xyzhotel/domain/customer"
	"xyzhotel/domain/money"
)

type CreditWalletCmd struct {
	CustomerID customer.ID
	Money      money.Money
}

type CreditWalletHandler struct {
	CustomerService *customer.Service
}

func (h CreditWalletHandler) Handle(ctx context.Context, cmd *CreditWalletCmd) error {
	cust, err := h.CustomerService.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return err
	}

	converter, err := money.CurrencyToConverter(cmd.Money.Currency)
	if err != nil {
		return err
	}

	convertedAmount := converter.Convert(cmd.Money)
	cust.Wallet.Credit(convertedAmount)
	return h.CustomerService.UpdateCustomer(ctx, cust)
}
