package application

import (
	"context"
	"testing"

	"xyzhotel/internal/domain/customer"

	"github.com/google/uuid"
)

func TestListCustomersHandler_Empty(t *testing.T) {
	repo := NewMockCustomerRepo()
	handler := &ListCustomersHandler{CustomerRepository: repo}

	cmd := &ListCustomersCmd{}
	list, err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(list) != 0 {
		t.Errorf("expected 0 customers, got %d", len(list))
	}
}

func TestListCustomersHandler_Success(t *testing.T) {
	repo := NewMockCustomerRepo()

	c1 := customer.NewCustomer(uuid.New(), "Alice", "alice@test.com", "111", customer.ZeroWallet())
	c2 := customer.NewCustomer(uuid.New(), "Bob", "bob@test.com", "222", customer.ZeroWallet())

	repo.CreateCustomer(context.Background(), c1)
	repo.CreateCustomer(context.Background(), c2)

	handler := &ListCustomersHandler{CustomerRepository: repo}

	cmd := &ListCustomersCmd{}
	list, err := handler.Handle(context.Background(), cmd)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 customers, got %d", len(list))
	}
}
