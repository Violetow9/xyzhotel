package customer

import (
	"context"
	"database/sql"
	"errors"

	"xyzhotel/domain/customer"
	"xyzhotel/infrastructure/mysql/sqldb"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound  = errors.New("customer not found")
	ErrDuplicate = errors.New("customer duplicate (id or email)")
)

type Repository struct {
	db *sql.DB
	q  *sqldb.Queries
}

func NewCustomerRepository(db *sql.DB) *Repository {
	return &Repository{db: db, q: sqldb.New(db)}
}

/* ===== CRUD ===== */

func (r *Repository) GetCustomerByID(ctx context.Context, id customer.ID) (*customer.Customer, error) {
	row, err := r.q.GetCustomerByID(ctx, id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toDomain(row.ID, row.Email, row.Phone, row.FullName, row.WalletSold)
}

func (r *Repository) GetCustomerByEmail(ctx context.Context, email customer.Email) (*customer.Customer, error) {
	row, err := r.q.GetCustomerByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return toDomain(row.ID, row.Email, row.Phone, row.FullName, row.WalletSold)
}

func (r *Repository) ListCustomers(ctx context.Context) ([]customer.Customer, error) {
	rows, err := r.q.ListCustomers(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]customer.Customer, 0, len(rows))
	for _, row := range rows {
		c, err := toDomain(row.ID, row.Email, row.Phone, row.FullName, row.WalletSold)
		if err != nil {
			return nil, err
		}
		out = append(out, *c)
	}
	return out, nil
}

func (r *Repository) CreateCustomer(ctx context.Context, c *customer.Customer) error {
	err := r.q.CreateCustomer(ctx, sqldb.CreateCustomerParams{
		ID:         c.ID.String(),
		Email:      c.Email,
		Phone:      c.Phone,
		FullName:   c.FullName,
		WalletSold: c.Wallet.Sold.String(),
	})
	if dup(err) {
		return ErrDuplicate
	}
	return err
}

func (r *Repository) UpdateCustomer(ctx context.Context, c *customer.Customer) error {
	aff, err := r.q.UpdateCustomer(ctx, sqldb.UpdateCustomerParams{
		Email:      c.Email,
		Phone:      c.Phone,
		FullName:   c.FullName,
		WalletSold: c.Wallet.Sold.String(),
		ID:         c.ID.String(),
	})
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteCustomer(ctx context.Context, id customer.ID) error {
	aff, err := r.q.DeleteCustomer(ctx, id.String())
	if err != nil {
		return err
	}
	if aff == 0 {
		return ErrNotFound
	}
	return nil
}

func dup(err error) bool {
	type myErr interface{ Number() uint16 }
	var me myErr
	return err != nil && errors.As(err, &me) && me.Number() == 1062
}

func toDomain(idStr, email, phone, fullName string, wallet string) (*customer.Customer, error) {
	uid, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	walletDec, err := decimal.NewFromString(wallet)
	if err != nil {
		return nil, err
	}

	w := customer.NewWallet(walletDec)
	return customer.NewCustomer(uid, fullName, email, phone, w), nil
}
