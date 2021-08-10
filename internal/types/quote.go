package types

const (
	CurrencyISO4217 = 3
)

type FullQuoteData struct {
	// ID is our own, not the external service's
	ID  string                  `json:"id"`
	CreateQuoteRequestData  `json:"req"`
	CreateQuoteResponseData `json:"res"`
}

type CreateQuoteRequestData struct {
	SourceCurrency string `json:"source_currency" db:"source_currency"`
	TargetCurrency string `json:"target_currency" db:"target_currency"`
	Amount         string `json:"amount" db:"amount"`
	AccountID      uint64 `json:"account_id" db:"account_id"`
}

type CreateQuoteResponseData struct {
	QuoteID        string `json:"quote_id" db:"quote_id"`
	TransactionFee string `json:"transaction_fee" db:"transaction_fee"`
	// EDT - estimated delivery time in seconds
	EDT int64 `json:"edt" db:"edt"`
}

type GetQuoteRequestData struct {
	ID        string `json:"id"`
	AccountID uint64 `json:"account_id"`
}
