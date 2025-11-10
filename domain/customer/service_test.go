package customer

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

var _ Repository = &InMemoryCustomerRepository{}

type InMemoryCustomerRepository struct {
	customersByEmail map[Email]*Customer
}

func (r InMemoryCustomerRepository) GetCustomerByID(_ context.Context, id ID) (*Customer, error) {
	for _, customer := range r.customersByEmail {
		if customer.ID == id {
			return customer, nil
		}
	}
	return nil, nil
}

func (r InMemoryCustomerRepository) GetCustomerByEmail(_ context.Context, email Email) (*Customer, error) {
	customer, exists := r.customersByEmail[email]
	if !exists {
		return nil, nil
	}
	return customer, nil
}

func (r InMemoryCustomerRepository) CreateCustomer(_ context.Context, customer *Customer) error {
	r.customersByEmail[customer.Email] = customer
	return nil
}

func (r InMemoryCustomerRepository) UpdateCustomer(_ context.Context, customer *Customer) error {
	r.customersByEmail[customer.Email] = customer
	return nil
}

func (r InMemoryCustomerRepository) DeleteCustomer(_ context.Context, id ID) error {
	for email, customer := range r.customersByEmail {
		if customer.ID == id {
			delete(r.customersByEmail, email)
			return nil
		}
	}
	return nil
}

func (r InMemoryCustomerRepository) ListCustomers(_ context.Context) ([]Customer, error) {
	customers := make([]Customer, 0, len(r.customersByEmail))
	for _, customer := range r.customersByEmail {
		customers = append(customers, *customer)
	}
	return customers, nil
}

func NewInMemoryCustomerRepository() *InMemoryCustomerRepository {
	return &InMemoryCustomerRepository{
		customersByEmail: make(map[Email]*Customer),
	}
}

func TestService_CreateCustomer_ExistingAccount(t *testing.T) {
	// Given a customer service
	// And an email that already exists in the Repository
	ctx := context.Background()
	repository := NewInMemoryCustomerRepository()
	email := Email("existingemail@gmail.com")
	err := repository.CreateCustomer(ctx, NewCustomer(uuid.New(), "Existing User", email, "1234567890", ZeroWallet()))
	if err != nil {
		t.Fatalf("failed to set up existing customer: %v", err)
	}
	customerService := Service{repository}

	// When creating a new customer with that email
	createCustomer := &CreateCommand{
		FullName: "NewCustomer",
		Email:    email,
		Phone:    "0987654321",
	}
	customer, err := customerService.CreateCustomer(ctx, createCustomer)
	if err == nil {
		t.Errorf("expected error but got none")
	}
	if customer != nil {
		t.Errorf("expected no customer to be created, but got one")
	}

	// Then expect an ErrCustomerAlreadyExists error
	if !errors.Is(err, ErrCustomerAlreadyExists) {
		t.Errorf("expected ErrCustomerAlreadyExists, but got: %v", err)
	}
}

func TestService_CreateCustomer_FieldEmpty(t *testing.T) {
	// Given a customer service
	// And a customer with an empty full name
	customerService := Service{NewInMemoryCustomerRepository()}

	// When creating a new customer with that email
	createCustomer := &CreateCommand{
		FullName: "",
		Email:    "",
		Phone:    "",
	}
	customer, err := customerService.CreateCustomer(context.Background(), createCustomer)
	if err == nil {
		t.Errorf("expected error but got none")
	}
	if customer != nil {
		t.Errorf("expected no customer to be created, but got one")
	}

	// Then expect an ErrCustomerFieldEmpty error
	if !errors.Is(err, ErrCustomerFieldEmpty) {
		t.Errorf("expected ErrCustomerFieldEmpty, but got: %v", err)
	}
}

func TestService_CreateCustomer_NonExistingAccount(t *testing.T) {
	// Given a customer service
	// And an email that does not exist in the Repository
	ctx := context.Background()
	repository := NewInMemoryCustomerRepository()
	err := repository.CreateCustomer(ctx, NewCustomer(uuid.New(), "Bob", "bob@gmail.com", "1234567890", ZeroWallet()))
	if err != nil {
		t.Fatalf("failed to set up existing customer: %v", err)
	}
	customerService := Service{repository}

	// When creating a new customer with that email
	email := "alice@gmail.com"
	createCustomer := &CreateCommand{
		FullName: "Alice",
		Email:    email,
		Phone:    "0987654321",
	}
	customer, err := customerService.CreateCustomer(ctx, createCustomer)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then expect the customer to be created successfully
	if customer == nil {
		t.Errorf("expected customer to be created, but got nil")
		return
	}

	if customer.Email != email {
		t.Errorf("expected customer email to be %s, but got %s", email, customer.Email)
	}

	if customer.FullName != "Alice" {
		t.Errorf("expected customer full name to be Alice, but got %s", customer.FullName)
	}

	if customer.Phone != "0987654321" {
		t.Errorf("expected customer phone to be 0987654321, but got %s", customer.Phone)
	}
}
