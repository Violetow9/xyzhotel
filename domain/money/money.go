package money

import "github.com/shopspring/decimal"

type Money struct {
	Amount   decimal.Decimal
	Currency Currency
}
