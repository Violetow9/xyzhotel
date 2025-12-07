package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"xyzhotel/internal/domain/customer"
	"xyzhotel/internal/infrastructure/mysql/sqldb"

	"github.com/google/uuid"
)

type CustomerRepository struct {
	db *sql.DB
	q  *sqldb.Queries
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
		q:  sqldb.New(db),
	}
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, c *customer.Customer) error {
	return r.q.CreateCustomer(ctx, sqldb.CreateCustomerParams{
		ID:           c.ID.String(),
		FullName:     c.FullName,
		Email:        c.Email,
		Phone:        c.Phone,
		BalanceCents: int32(c.Wallet.Balance()),
	})
}

func (r *CustomerRepository) GetCustomerByID(ctx context.Context, id customer.ID) (*customer.Customer, error) {
	row, err := r.q.GetCustomerByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		return nil, err
	}

	return r.mapToDomain(row)
}

func (r *CustomerRepository) GetCustomerByEmail(ctx context.Context, email string) (*customer.Customer, error) {
	row, err := r.q.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.mapToDomain(row)
}

func (r *CustomerRepository) UpdateCustomer(ctx context.Context, c *customer.Customer) error {
	return r.q.UpdateCustomer(ctx, sqldb.UpdateCustomerParams{
		FullName:     c.FullName,
		Email:        c.Email,
		Phone:        c.Phone,
		BalanceCents: int32(c.Wallet.Balance()),
		ID:           c.ID.String(),
	})
}

func (r *CustomerRepository) ListCustomers(ctx context.Context) ([]customer.Customer, error) {
	rows, err := r.q.ListCustomers(ctx)
	if err != nil {
		return nil, err
	}

	customers := make([]customer.Customer, 0, len(rows))
	for _, row := range rows {
		c, err := r.mapToDomain(row)
		if err != nil {
			return nil, err
		}
		customers = append(customers, *c)
	}
	return customers, nil
}

func (r *CustomerRepository) mapToDomain(row sqldb.Customer) (*customer.Customer, error) {
	id, err := uuid.Parse(row.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid uuid in db: %w", err)
	}

	wallet := customer.NewWallet(int(row.BalanceCents))

	return customer.NewCustomer(
		id,
		row.FullName,
		row.Email,
		row.Phone,
		wallet,
	), nil
}

func (r *CustomerRepository) DeleteCustomer(ctx context.Context, id customer.ID) error {
	return errors.New("not implemented")
}
