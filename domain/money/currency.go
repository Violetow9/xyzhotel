package money

import "fmt"

type Currency = string

const (
	EUR Currency = "EUR"
	USD Currency = "USD"
	LS  Currency = "LS"
	YEN Currency = "YEN"
	FS  Currency = "FS"
)

var ErrInvalidCurrency = fmt.Errorf("invalid currency")

func CurrencyFromString(s string) (Currency, error) {
	switch s {
	case "USD":
		return USD, nil
	case "EUR":
		return EUR, nil
	case "LS":
		return LS, nil
	case "YEN":
		return YEN, nil
	case "FS":
		return FS, nil
	default:
		return "", ErrInvalidCurrency
	}
}
