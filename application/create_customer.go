package application

import (
	"context"
	"errors"

	"xyzhotel/domain/customer"

	"github.com/google/uuid"
)

type CreateCustomerCmd struct {
	FullName string
	Email    customer.Email
	Phone    customer.Email
}

func (cc *CreateCustomerCmd) IsEmpty() bool {
	return cc.FullName == "" || cc.Email == "" || cc.Phone == ""
}

var (
	ErrCustomerAlreadyExists = errors.New("customer already exists")
	ErrCustomerFieldEmpty    = errors.New("customer field is empty")
)

type CreateCustomerHandler struct {
	CustomerRepository customer.Repository
}

func (h CreateCustomerHandler) Handle(ctx context.Context, cmd *CreateCustomerCmd) (*customer.Customer, error) {
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

	cust = &customer.Customer{
		ID:       id,
		FullName: cmd.FullName,
		Email:    cmd.Email,
		Phone:    cmd.Phone,
		Wallet:   customer.ZeroWallet(),
	}
	err = h.CustomerRepository.CreateCustomer(ctx, cust)
	if err != nil {
		return nil, err
	}
	return cust, nil
}
