package money

import "fmt"

type Currency string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
	GBP Currency = "GBP"
	JPY Currency = "JPY"
	CHF Currency = "CHF"
)

type Money struct {
	AmountCents int
	Currency    Currency
}

func New(amount int, currency Currency) Money {
	return Money{
		AmountCents: amount,
		Currency:    currency,
	}
}

func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", float64(m.AmountCents)/100.0, m.Currency)
}
