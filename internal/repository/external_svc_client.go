package repository

import (
	"context"
	"github.com/pkg/errors"
	"some-http-server/internal/types"
)

type XPConfig struct {
	// Here we could have some webhook address, baseURL for requests, etc
}

type ExClientStub struct {
	xpAPI XPConfig
}

func NewExClientStub(xpAPI XPConfig) *ExClientStub {
	return &ExClientStub{xpAPI: xpAPI}
}

type GetQuoteRequest struct {
	Data types.QuoteRequest
}

type GetQuoteResponse struct {
	Data types.QuoteResponse
}

// GetQuote of the ExClientStub currently returns hard-coded data
func (c *ExClientStub) GetQuote(ctx context.Context, data *types.QuoteRequest) (*GetQuoteResponse, error) {
	if data == nil {
		return nil, errors.New("empty quote request data")
	}

	// Creating request here
	req := GetQuoteRequest{Data: *data}

	_ = ctx

	return &GetQuoteResponse{
		Data: types.QuoteResponse{
			ID:             "some_id",
			TransactionFee: float64(req.Data.Amount) * 0.07,
			EDT:            120,
		},
	}, nil
}
