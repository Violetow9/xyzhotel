package application

import (
	"context"
	"errors"
	"testing"

	"xyzhotel/internal/domain/customer"

	"github.com/google/uuid"
)

func TestCreateCustomerHandler_EmptyFields(t *testing.T) {
	handler := &CreateCustomerHandler{}

	cmd := &CreateCustomerCmd{
		FullName: "",
		Email:    "",
		Phone:    "",
	}

	_, err := handler.Handle(context.Background(), cmd)

	if err == nil {
		t.Error("expected error for empty fields")
	}
}

func TestCreateCustomerHandler_CustomerAlreadyExists(t *testing.T) {
	repo := NewMockCustomerRepo()
	existingEmail := "alice@example.com"

	existingCustomer := customer.NewCustomer(
		uuid.New(),
		"Alice",
		existingEmail,
		"123456789",
		customer.ZeroWallet(),
	)
	repo.CreateCustomer(context.Background(), existingCustomer)

	handler := &CreateCustomerHandler{CustomerRepository: repo}

	cmd := &CreateCustomerCmd{
		FullName: "Alice Duplicate",
		Email:    existingEmail,
		Phone:    "987654321",
	}

	_, err := handler.Handle(context.Background(), cmd)

	if !errors.Is(err, ErrCustomerAlreadyExists) {
		t.Errorf("expected ErrCustomerAlreadyExists, got %v", err)
	}
}

func TestCreateCustomerHandler_Success(t *testing.T) {
	repo := NewMockCustomerRepo()
	handler := &CreateCustomerHandler{CustomerRepository: repo}

	cmd := &CreateCustomerCmd{
		FullName: "Bob",
		Email:    "bob@example.com",
		Phone:    "123456789",
	}

	cust, err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cust.ID == uuid.Nil {
		t.Error("expected valid UUID")
	}

	if cust.Wallet.Balance() != 0 {
		t.Error("new customer should have zero balance")
	}

	savedCust, _ := repo.GetCustomerByID(context.Background(), cust.ID)
	if savedCust == nil {
		t.Error("customer was not saved in repository")
	}
}
