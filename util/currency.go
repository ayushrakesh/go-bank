package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(cur string) bool {
	switch cur {
	case USD, EUR, CAD:
		return true
	}
	return false
}
