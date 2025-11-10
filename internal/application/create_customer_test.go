package application

import (
	"context"
	"errors"
	"testing"
	customer2 "xyzhotel/internal/domain/customer"

	"github.com/google/uuid"
)

func TestCreateCustomerHandler_EmailAlreadyExists(t *testing.T) {
	// Given a create customer handler
	// And an email that already exists in the Repository
	ctx := context.Background()
	email := "existingemail@gmail.com"
	repository := customer2.NewFakeRepositoryWithCustomers(&customer2.Customer{
		ID:       uuid.New(),
		FullName: "ExistingUser",
		Email:    email,
		Phone:    "1234567890",
		Wallet:   customer2.ZeroWallet(),
	})
	createCustomerHandler := &CreateCustomerHandler{
		CustomerRepository: repository,
	}

	// When creating a new customer with that email
	cust, err := createCustomerHandler.Handle(ctx, &CreateCustomerCmd{
		FullName: "NewCustomer",
		Email:    email,
		Phone:    "0987654321",
	})
	if err == nil {
		t.Errorf("expected error but got none")
	}
	if cust != nil {
		t.Errorf("expected no customer to be created, but got one")
	}

	// Then expect an ErrCustomerAlreadyExists error
	if !errors.Is(err, ErrCustomerAlreadyExists) {
		t.Errorf("expected ErrCustomerAlreadyExists, but got: %v", err)
	}
}

func TestCreateCustomerHandler_FieldEmpty(t *testing.T) {
	// Given a create customer handler
	// And a customer with an empty full name
	ctx := context.Background()
	repository := customer2.NewFakeRepository()
	createCustomerHandler := &CreateCustomerHandler{
		CustomerRepository: repository,
	}

	// When creating a new customer with empty fields
	cust, err := createCustomerHandler.Handle(ctx, &CreateCustomerCmd{
		FullName: "",
		Email:    "",
		Phone:    "",
	})
	if err == nil {
		t.Errorf("expected error but got none")
	}
	if cust != nil {
		t.Errorf("expected no customer to be created, but got one")
	}

	// Then expect an ErrCustomerFieldEmpty error
	if !errors.Is(err, ErrCustomerFieldEmpty) {
		t.Errorf("expected ErrCustomerFieldEmpty, but got: %v", err)
	}
}

func TestCreateCustomerHandler_NonExistingAccount(t *testing.T) {
	// Given a customer service
	// And an email that does not exist in the Repository
	ctx := context.Background()
	repository := customer2.NewFakeRepository()
	createCustomerHandler := &CreateCustomerHandler{
		CustomerRepository: repository,
	}

	// When creating a new customer
	cust, err := createCustomerHandler.Handle(ctx, &CreateCustomerCmd{
		FullName: "Alice",
		Email:    "alice@gmail.com",
		Phone:    "1234567890",
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Then expect the customer to be created successfully
	if cust == nil {
		t.Errorf("expected customer to be created, but got nil")
		return
	}

	if cust.Email != "alice@gmail.com" {
		t.Errorf("expected customer email to be alice@gmail.com, but got %s", cust.Email)
	}

	if cust.FullName != "Alice" {
		t.Errorf("expected customer full name to be Alice, but got %s", cust.FullName)
	}

	if cust.Phone != "1234567890" {
		t.Errorf("expected customer phone to be 1234567890, but got %s", cust.Phone)
	}
}
