package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"some-http-server/internal/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

type ExternalSvcClient interface {
	CreateQuote(ctx context.Context, data *types.CreateQuoteRequestData) (*types.CreateQuoteResponseData, error)
}

type QuoteRepo interface {
	Save(ctx context.Context, quote *types.FullQuoteData) (string, error)
	Read(ctx context.Context, id string) (*types.FullQuoteData, error)
}

type QuoteService struct {
	xSvc      ExternalSvcClient
	quoteRepo QuoteRepo
}

func NewQuoteService(xSvc ExternalSvcClient, quoteRepo QuoteRepo) *QuoteService {
	return &QuoteService{xSvc: xSvc, quoteRepo: quoteRepo}
}

func (s *QuoteService) Create(ctx context.Context, req *types.CreateQuoteRequestData) (*types.CreateQuoteResponseData, error) {
	if req == nil {
		return nil, types.ErrEmptyRequest
	}
	log := ctxlogrus.Extract(ctx)

	res, err := s.xSvc.CreateQuote(ctx, req)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create a quote")
	}

	log.Info("Successfully created a quote")

	fqd := &types.FullQuoteData{Req: req, Res: res}

	id, err := s.quoteRepo.Save(ctx, fqd)
	if err != nil {
		return nil, errors.Wrap(err, "cannot save the quote")
	}
	log.Info("Successfully saved the quote in the db")

	// Not to expose real ID and avoid relying on external svc provider in its consistency
	res.QuoteID = id

	return res, nil
}

func (s *QuoteService) Read(ctx context.Context, req *types.GetQuoteRequestData) (*types.FullQuoteData, error) {
	if req == nil {
		return nil, types.ErrEmptyRequest
	}
	log := ctxlogrus.Extract(ctx)
	log.Infof("Trying to read the quote %s", req.ID)

	fqd, err := s.quoteRepo.Read(ctx, fmt.Sprint(req.ID))
	if err != nil {
		return &types.FullQuoteData{}, errors.Wrap(err, "cannot read the quote")
	}

	log.Infof("Successfully read the quote %s", req.ID)

	return fqd, nil
}
