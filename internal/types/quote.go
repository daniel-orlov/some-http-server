package types

type QuoteRequest struct {
	SourceCurrency string `json:"source_currency" db:"source_currency"`
	TargetCurrency string `json:"target_currency" db:"target_currency"`
	Amount    int64 `json:"amount" db:"amount"`
	AccountID uint64 `json:"account_id" db:"account_id"`
}

type QuoteResponse struct {
	ID             string  `json:"id" db:"quote_id"`
	TransactionFee float64 `json:"transaction_fee" db:"transaction_fee"`
	// EDT - estimated delivery time
	EDT int64 `json:"edt" db:"edt"`
}

type FullQuoteData struct {
	// ID is our own, not the external service's
	ID  string        `json:"id"`
	Req QuoteRequest  `json:"req"`
	Res QuoteResponse `json:"res"`
}
