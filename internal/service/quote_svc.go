package service

import (
	"context"
	"github.com/pkg/errors"
	repo "some-http-server/internal/repository"
	"some-http-server/internal/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

type ExternalSvcClient interface {
	GetQuote(ctx context.Context, req *repo.GetQuoteRequest) (*repo.GetQuoteResponse, error)
}

//github.com/myntra/golimit

type LocalQuoteRepo interface {
	Save(ctx context.Context, quote *types.FullQuoteData) (string, error)
	Read(ctx context.Context, id string) (*types.FullQuoteData, error)
}

type QuoteService struct {
	xSvc  ExternalSvcClient
	local LocalQuoteRepo
}

func NewQuoteService(xSvc ExternalSvcClient, local LocalQuoteRepo) *QuoteService {
	return &QuoteService{xSvc: xSvc, local: local}
}

func (s *QuoteService) Create(ctx context.Context, q *types.QuoteRequest) (string, error) {
	if q == nil {
		return "", errors.New("empty quote request data")
	}
	log := ctxlogrus.Extract(ctx)

	qd, err := s.xSvc.GetQuote(ctx, q)
	if err != nil {
		return "", errors.Wrap(err, "cannot create a quote")
	}

	log.Info("successfully created a quote")

	fqd := types.FullQuoteData{Req: quote, Res: qd}

	id, err := s.local.Save(ctx, fqd)
	if err != nil {
		return "", errors.Wrap(err, "cannot save the quote")
	}

	log.Info("successfully saved the quote")

	return id, nil
}

func (s *QuoteService) Read(ctx context.Context, id string) (types.FullQuoteData, error) {
	log := ctxlogrus.Extract(ctx)
	log.Infof("trying to read the quote %s", id)

	fqd, err := s.local.Read(ctx, id)
	if err != nil {
		return types.FullQuoteData{}, errors.Wrap(err, "cannot read the quote")
	}

	log.Infof("successfully read the quote %s", id)

	return fqd, nil
}
