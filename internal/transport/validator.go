package transport

import (
	"some-http-server/internal/types"
	"strconv"
)

func ValidateCreateQuoteRequest(req *CreateQuoteRequest) error {
	if !validateCurrencyIsISO4217(req.Data.SourceCurrency) || !validateCurrencyIsISO4217(req.Data.TargetCurrency) {
		return types.ErrCurrencyNotISO4217
	}

	if !validateStringParsesToFloatAndNotZero(req.Data.Amount) {
		return types.ErrAmountInvalid
	}

	return nil
}

func validateStringParsesToFloatAndNotZero(s string) bool {
	result, err := strconv.ParseFloat(s, 64)
	return err == nil && result != float64(0)
}

func validateCurrencyIsISO4217(s string) bool {
	if !validateStringIsLen(s, types.CurrencyISO4217) {
		return false
	}

	for i := range s {
		if (s[i] < 'a' || s[i] > 'z') && (s[i] < 'A' || s[i] > 'Z') {
			return false
		}
	}

	return true
}

func validateStringIsLen(s string, l int) bool {
	return len(s) == l
}
