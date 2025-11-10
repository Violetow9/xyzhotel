package money

import "errors"

func CurrencyToConverter(currency Currency) (*Converter, error) {
	var converter *Converter
	switch currency {
	case USD:
		converter = USDToEURConverter
	case EUR:
		converter = EURToEURConverter
	case LS:
		converter = LSToEURConverter
	case YEN:
		converter = YENToEURConverter
	case FS:
		converter = FSToEURConverter
	default:
		return nil, errors.New("unsupported currency")
	}
	return converter, nil
}
