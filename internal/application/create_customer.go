package application

import (
	"context"
	"errors"
	customer2 "xyzhotel/internal/domain/customer"

	"github.com/google/uuid"
)

type CreateCustomerCmd struct {
	FullName string
	Email    customer2.Email
	Phone    customer2.Email
}

func (cc *CreateCustomerCmd) IsEmpty() bool {
	return cc.FullName == "" || cc.Email == "" || cc.Phone == ""
}

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrCustomerFieldEmpty    = errors.New("customer field is empty")
)

type CreateCustomerHandler struct {
	CustomerRepository customer2.Repository
}

func (h CreateCustomerHandler) Handle(ctx context.Context, cmd *CreateCustomerCmd) (*customer2.Customer, error) {
	if cmd.IsEmpty() {
		return nil, ErrCustomerFieldEmpty
	}

	cust, err := h.CustomerRepository.GetCustomerByEmail(ctx, cmd.Email)
	if err == nil && cust != nil {
		return nil, ErrCustomerAlreadyExists
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	cust = &customer2.Customer{
		ID:       id,
		FullName: cmd.FullName,
		Email:    cmd.Email,
		Phone:    cmd.Phone,
		Wallet:   customer2.ZeroWallet(),
	}
	err = h.CustomerRepository.CreateCustomer(ctx, cust)
	if err != nil {
		return nil, err
	}
	return cust, nil
}
