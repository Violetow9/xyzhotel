package customer

import (
	"context"
	"errors"
)

var ErrCustomerAlreadyExists = errors.New("customer already exists")
var ErrCustomerNotFound = errors.New("customer not found")
var ErrCustomerFieldEmpty = errors.New("customer field is empty")

type Service struct {
	Repository Repository
}

func (cs *Service) GetCustomerByID(ctx context.Context, id ID) (*Customer, error) {
	return cs.Repository.GetCustomerByID(ctx, id)
}

func (cs *Service) UpdateCustomer(ctx context.Context, customer *Customer) error {
	return cs.Repository.UpdateCustomer(ctx, customer)
}

func (cs *Service) DeleteCustomer(ctx context.Context, id ID) error {
	return cs.Repository.DeleteCustomer(ctx, id)
}

func (cs *Service) ListCustomers(ctx context.Context) ([]Customer, error) {
	return cs.Repository.ListCustomers(ctx)
}
