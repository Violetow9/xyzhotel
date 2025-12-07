package application

import (
	"context"
	"errors"
	"fmt"
	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/domain/money"
)

type DebitWalletCmd struct {
	CustomerID customer.ID
	Money      money.Money
}

var (
	ErrInsufficientFunds = errors.New("not sufficient funds")
)

type DebitWalletHandler struct {
	CustomerRepo customer.Repository
	Converter    *money.Converter
}

func (h DebitWalletHandler) Handle(ctx context.Context, cmd *DebitWalletCmd) error {
	cust, err := h.CustomerRepo.GetCustomerByID(ctx, cmd.CustomerID)
	if err != nil {
		return fmt.Errorf("failed to fetch customer: %w", err)
	}

	amountInEuros, err := h.Converter.ConvertToEUR(cmd.Money)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	if !cust.Wallet.HasSufficientFunds(amountInEuros.AmountCents) {
		return ErrInsufficientFunds
	}

	if err := cust.Wallet.Debit(amountInEuros.AmountCents); err != nil {
		return err
	}

	if err := h.CustomerRepo.UpdateCustomer(ctx, cust); err != nil {
		return fmt.Errorf("failed to update customer: %w", err)
	}

	return nil
}
