package types

import "github.com/pkg/errors"

var (
	ErrEmptyRequest = errors.New("empty request data")

	ErrCurrencyNotISO4217 = errors.New("currency violates ISO-4217 standard")
	ErrAmountInvalid      = errors.New("amount is invalid")
)
