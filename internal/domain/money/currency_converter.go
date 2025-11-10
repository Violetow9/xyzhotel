package money

import (
	"github.com/shopspring/decimal"
)

type Converter struct {
	FromDevise Currency
	ToDevise   Currency
	Rate       decimal.Decimal
}

func NewDeviseConverter(from Currency, to Currency, rate decimal.Decimal) *Converter {
	return &Converter{
		FromDevise: from,
		ToDevise:   to,
		Rate:       rate,
	}
}

func (wc *Converter) Convert(money Money) decimal.Decimal {
	return money.Amount.Mul(wc.Rate)
}

var (
	USDToEURRate      = decimal.NewFromFloat(0.85)
	USDToEURConverter = NewDeviseConverter(USD, EUR, USDToEURRate)
	EURToEURConverter = NewDeviseConverter(EUR, EUR, decimal.NewFromInt(1))
	LSToEURConverter  = NewDeviseConverter(LS, EUR, decimal.NewFromFloat(1.17))
	YENToEURConverter = NewDeviseConverter(YEN, EUR, decimal.NewFromFloat(0.0077))
	FSToEURConverter  = NewDeviseConverter(FS, EUR, decimal.NewFromFloat(0.92))
)
