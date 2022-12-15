package util

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	CAD = "CAD"
	JPY = "JPY"
	TRY = "TRY"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, CAD, JPY, TRY:
		return true
	default:
		return false
	}
}
