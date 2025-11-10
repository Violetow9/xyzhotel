package customer

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ID = uuid.UUID
type Email = string
type Phone = string

type Wallet struct{ Sold decimal.Decimal }

func NewWallet(initial decimal.Decimal) *Wallet             { return &Wallet{Sold: initial} }
func ZeroWallet() *Wallet                                   { return NewWallet(decimal.Zero) }
func (w *Wallet) Balance() decimal.Decimal                  { return w.Sold }
func (w *Wallet) Credit(amount decimal.Decimal)             { w.Sold = w.Sold.Add(amount) }
func (w *Wallet) Debit(amount decimal.Decimal)              { w.Sold = w.Sold.Sub(amount) }
func (w *Wallet) HasSufficientFunds(a decimal.Decimal) bool { return w.Sold.GreaterThanOrEqual(a) }

type Customer struct {
	ID       ID
	FullName string
	Email    Email
	Phone    Phone
	Wallet   *Wallet
}

func NewCustomer(id ID, fullName string, email Email, phone Phone, wallet *Wallet) *Customer {
	return &Customer{ID: id, FullName: fullName, Email: email, Phone: phone, Wallet: wallet}
}
