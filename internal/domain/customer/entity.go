package customer

import (
	"errors"

	"github.com/google/uuid"
)

type ID = uuid.UUID
type Email = string
type Phone = string

var (
	ErrInsufficientFunds = errors.New("insufficient funds in wallet")
)

type Wallet struct {
	BalanceCents int
}

func NewWallet(initialCents int) *Wallet {
	return &Wallet{BalanceCents: initialCents}
}

func ZeroWallet() *Wallet {
	return NewWallet(0)
}

func (w *Wallet) Balance() int {
	return w.BalanceCents
}

func (w *Wallet) Credit(amountCents int) {
	if amountCents < 0 {
		return
	}
	w.BalanceCents += amountCents
}

func (w *Wallet) Debit(amountCents int) error {
	if amountCents < 0 {
		return errors.New("cannot debit negative amount")
	}
	if !w.HasSufficientFunds(amountCents) {
		return ErrInsufficientFunds
	}
	w.BalanceCents -= amountCents
	return nil
}

func (w *Wallet) HasSufficientFunds(amountCents int) bool {
	return w.BalanceCents >= amountCents
}

type Customer struct {
	ID       ID
	FullName string
	Email    Email
	Phone    Phone
	Wallet   *Wallet
}

func NewCustomer(id ID, fullName string, email Email, phone Phone, wallet *Wallet) *Customer {
	return &Customer{
		ID:       id,
		FullName: fullName,
		Email:    email,
		Phone:    phone,
		Wallet:   wallet,
	}
}
