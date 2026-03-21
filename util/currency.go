package util

var supportedCurrencies = map[string]struct{}{
	"USD": {},
	"CAD": {},
	"EUR": {},
	"INR": {},
}

func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencies[currency]
	return ok
}