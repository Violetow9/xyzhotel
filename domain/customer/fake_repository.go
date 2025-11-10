package customer

import "context"

var _ Repository = &FakeRepository{}

type FakeRepository struct {
	Customers map[ID]*Customer
}

func (r *FakeRepository) GetCustomerByID(ctx context.Context, id ID) (*Customer, error) {
	customer, exists := r.Customers[id]
	if !exists {
		return nil, ErrCustomerNotFound
	}
	return customer, nil
}

func (r *FakeRepository) GetCustomerByEmail(ctx context.Context, email Email) (*Customer, error) {
	for _, customer := range r.Customers {
		if customer.Email == email {
			return customer, nil
		}
	}
	return nil, nil
}

func (r *FakeRepository) CreateCustomer(ctx context.Context, customer *Customer) error {
	r.Customers[customer.ID] = customer
	return nil
}

func (r *FakeRepository) UpdateCustomer(ctx context.Context, customer *Customer) error {
	r.Customers[customer.ID] = customer
	return nil
}

func (r *FakeRepository) DeleteCustomer(ctx context.Context, id ID) error {
	delete(r.Customers, id)
	return nil
}

func (r *FakeRepository) ListCustomers(ctx context.Context) ([]Customer, error) {
	customers := make([]Customer, 0, len(r.Customers))
	for _, customer := range r.Customers {
		customers = append(customers, *customer)
	}
	return customers, nil
}

func NewFakeRepository() *FakeRepository {
	return &FakeRepository{
		Customers: make(map[ID]*Customer),
	}
}

func NewFakeRepositoryWithCustomers(customers ...*Customer) *FakeRepository {
	repo := &FakeRepository{
		Customers: make(map[ID]*Customer),
	}
	for _, customer := range customers {
		repo.Customers[customer.ID] = customer
	}
	return repo
}
