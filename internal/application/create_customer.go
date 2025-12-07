package application

import (
	"context"
	"errors"
	"fmt"

	"xyzhotel/internal/domain/customer"

	"github.com/google/uuid"
)

type CreateCustomerCmd struct {
	FullName string
	Email    customer.Email
	Phone    customer.Phone
}

func (cc *CreateCustomerCmd) Validate() error {
	if cc.FullName == "" || cc.Email == "" || cc.Phone == "" {
		return errors.New("customer fields cannot be empty")
	}
	return nil
}

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
)

type CreateCustomerHandler struct {
	CustomerRepository customer.Repository
}

func (h CreateCustomerHandler) Handle(ctx context.Context, cmd *CreateCustomerCmd) (*customer.Customer, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	existing, err := h.CustomerRepository.GetCustomerByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check customer existence: %w", err)
	}
	if existing != nil {
		return nil, ErrCustomerAlreadyExists
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	cust := customer.NewCustomer(
		id,
		cmd.FullName,
		cmd.Email,
		cmd.Phone,
		customer.ZeroWallet(),
	)

	if err := h.CustomerRepository.CreateCustomer(ctx, cust); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return cust, nil
}
