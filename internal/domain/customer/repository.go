package customer

import "context"

type Repository interface {
	GetCustomerByID(ctx context.Context, id ID) (*Customer, error)

	GetCustomerByEmail(ctx context.Context, email Email) (*Customer, error)

	CreateCustomer(ctx context.Context, customer *Customer) error

	UpdateCustomer(ctx context.Context, customer *Customer) error

	DeleteCustomer(ctx context.Context, id ID) error

	ListCustomers(ctx context.Context) ([]Customer, error)
}
