package application

import (
	"context"
	"xyzhotel/internal/domain/customer"
)

type ListCustomersCmd struct {
}

type ListCustomersHandler struct {
	CustomerRepository customer.Repository
}

func (h ListCustomersHandler) Handle(ctx context.Context, cmd *ListCustomersCmd) ([]customer.Customer, error) {
	return h.CustomerRepository.ListCustomers(ctx)
}
