package utils

const (
	AUD = "AUD"
	GBP = "GBP"
	HKD = "HKD"
	IDR = "IDR"
	MYR = "MYR"
	NZD = "NZD"
	PHP = "PHP"
	SGD = "SGD"
	THB = "THB"
	VND = "VND"
)

// isSupportedCurrency returns true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case AUD, GBP, HKD, IDR, MYR, NZD, PHP, SGD, THB, VND:
		return true
	}
	return false
}
