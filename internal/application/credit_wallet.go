package application

import (
	"context"
	customer2 "xyzhotel/internal/domain/customer"
	money2 "xyzhotel/internal/domain/money"
)

type CreditWalletCmd struct {
	CustomerID customer2.ID
	Money      money2.Money
}

type CreditWalletHandler struct {
	CustomerService *customer2.Service
}

func (h CreditWalletHandler) Handle(ctx context.Context, cmd *CreditWalletCmd) error {
	cust, err := h.CustomerService.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return err
	}

	converter, err := money2.CurrencyToConverter(cmd.Money.Currency)
	if err != nil {
		return err
	}

	convertedAmount := converter.Convert(cmd.Money)
	cust.Wallet.Credit(convertedAmount)
	return h.CustomerService.UpdateCustomer(ctx, cust)
}
