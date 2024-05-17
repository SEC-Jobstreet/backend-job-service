package utils

const (
	REVIEW = "REVIEW"
	POSTED = "POSTED"
	DENIED = "DENIED"
	CLOSED = "CLOSED"
)

// isSupportedCurrency returns true if the currency is supported
func IsSupportedStatus(status string) bool {
	switch status {
	case REVIEW, POSTED, DENIED, CLOSED:
		return true
	}
	return false
}
