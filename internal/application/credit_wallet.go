package application

import (
	"context"
	"fmt"
	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"
)

type CreditWalletCmd struct {
	CustomerID customer.ID
	Money      money.Money
}

type CreditWalletHandler struct {
	CustomerRepo customer.Repository // On utilise le Repo, pas le Service (si possible)
	Converter    *money.Converter
}

func (h CreditWalletHandler) Handle(ctx context.Context, cmd *CreditWalletCmd) error {
	cust, err := h.CustomerRepo.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to fetch customer: %w", err)
	}

	amountInEuros, err := h.Converter.ConvertToEUR(cmd.Money)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	cust.Wallet.Credit(amountInEuros.AmountCents)

	if err := h.CustomerRepo.UpdateCustomer(ctx, cust); err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}
