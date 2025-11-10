package application

import (
	"context"
	customer2 "xyzhotel/internal/domain/customer"
)

type ListCustomersCmd struct {
}

type ListCustomersHandler struct {
	CustomerRepository customer2.Repository
}

func (h ListCustomersHandler) Handle(ctx context.Context, cmd *ListCustomersCmd) ([]customer2.Customer, error) {
	return h.CustomerRepository.ListCustomers(ctx)
}
