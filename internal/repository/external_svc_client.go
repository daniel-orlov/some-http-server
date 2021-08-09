package repository

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/pkg/errors"
	"some-http-server/internal/types"
)

type XPConfig struct {
	// Here we could have some webhook address, baseURL for requests, etc
}

type ExternalSvcClientStub struct {
	xpAPI XPConfig
}

func NewExternalSvcClientStub(xpAPI XPConfig) *ExternalSvcClientStub {
	return &ExternalSvcClientStub{xpAPI: xpAPI}
}

type CreateQuoteRequest struct {
	Data *types.CreateQuoteRequestData
}

type CreateQuoteResponse struct {
	Data *types.CreateQuoteResponseData
}

// CreateQuote of the ExClientStub currently returns hard-coded data
func (c *ExternalSvcClientStub) CreateQuote(ctx context.Context, data *types.CreateQuoteRequestData) (*types.CreateQuoteResponseData, error) {
	if data == nil {
		return nil, errors.New("empty quote request data")
	}
	log := ctxlogrus.Extract(ctx)

	// Creating request here
	req := CreateQuoteRequest{Data: data}
	log.Tracef("Sending a request '%v'", req)
	// Here a call to external service takes place

	fakeRes := &CreateQuoteResponse{
		Data: &types.CreateQuoteResponseData{
			QuoteID:        "some_id",
			TransactionFee: fmt.Sprint(123 * 0.07),
			EDT:            120,
		},
	}
	log.Tracef("Received a response '%v'", fakeRes)

	return fakeRes.GetData(), nil
}

// GetData safely extracts Data from GetQuoteResponse
func (r *CreateQuoteResponse) GetData() *types.CreateQuoteResponseData {
	if r == nil || r.Data == nil {
		return &types.CreateQuoteResponseData{}
	}

	return r.Data
}
