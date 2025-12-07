package money

import (
	"errors"
	"math"
)

var ErrUnsupportedCurrency = errors.New("unsupported currency")

type Converter struct {
	rates map[Currency]float64
}

func NewConverter() *Converter {
	return &Converter{
		rates: map[Currency]float64{
			EUR: 1.0,
			USD: 0.85,   // 1 USD = 0.85 EUR
			GBP: 1.17,   // 1 GBP = 1.17 EUR
			JPY: 0.0077, // 1 JPY = 0.0077 EUR
			CHF: 0.92,   // 1 CHF = 0.92 EUR
		},
	}
}

func (c *Converter) ConvertToEUR(money Money) (Money, error) {
	rate, exists := c.rates[money.Currency]
	if !exists {
		return Money{}, ErrUnsupportedCurrency
	}

	convertedAmount := math.Round(float64(money.AmountCents) * rate)

	return Money{
		AmountCents: int(convertedAmount),
		Currency:    EUR,
	}, nil
}
