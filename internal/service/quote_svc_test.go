package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"some-http-server/internal/service/mock"
	"some-http-server/internal/types"
	"testing"

	"github.com/golang/mock/gomock"
)

//go:generate mockgen -source quote_svc.go -package=mock -destination=mock/quote_svc_mock.go

func TestQuoteService_Create(t *testing.T) {
	ctx := context.Background()
	quoteID := "1234"

	request := &types.CreateQuoteRequestData{
		SourceCurrency: "BYN",
		TargetCurrency: "EUR",
		Amount:         "2021",
		AccountID:      42,
	}

	res := &types.CreateQuoteResponseData{
		QuoteID:        "quote_id_from_external_provider",
		TransactionFee: fmt.Sprint(123 * 0.07),
		EDT:            120,
	}

	fqd := &types.FullQuoteData{Req: request, Res: res}

	tests := []struct {
		name    string
		req     *types.CreateQuoteRequestData
		prepare func(xs *mock.MockExternalSvcClient, qr *mock.MockQuoteRepo)
		want    *types.CreateQuoteResponseData
		wantErr bool
	}{
		{
			"1. Error empty request",
			nil,
			func(xs *mock.MockExternalSvcClient, qr *mock.MockQuoteRepo) {},
			nil,
			true,
		},
		{
			"2. Error on create quote",
			request,
			func(xs *mock.MockExternalSvcClient, qr *mock.MockQuoteRepo) {
				xs.EXPECT().CreateQuote(ctx, request).Return(nil, errors.New("some error"))
			},
			nil,
			true,
		},
		{
			"3. Error on save quote",
			request,
			func(xs *mock.MockExternalSvcClient, qr *mock.MockQuoteRepo) {
				xs.EXPECT().CreateQuote(ctx, request).Return(res, nil)
				qr.EXPECT().Save(ctx, fqd).Return("", errors.New("some error"))
			},
			nil,
			true,
		},
		{
			"4. Success",
			request,
			func(xs *mock.MockExternalSvcClient, qr *mock.MockQuoteRepo) {
				xs.EXPECT().CreateQuote(ctx, request).Return(res, nil)
				qr.EXPECT().Save(ctx, fqd).Return(quoteID, nil)
			},
			&types.CreateQuoteResponseData{
				QuoteID:        quoteID,
				TransactionFee: fmt.Sprint(123 * 0.07),
				EDT:            120,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			xSvcMock := mock.NewMockExternalSvcClient(ctrl)
			quoteRepoMock := mock.NewMockQuoteRepo(ctrl)

			tt.prepare(xSvcMock, quoteRepoMock)

			s := &QuoteService{
				xSvc:      xSvcMock,
				quoteRepo: quoteRepoMock,
			}
			got, err := s.Create(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuoteService_Read(t *testing.T) {
	ctx := context.Background()
	quoteID := uint64(1234)
	accountID := uint64(42)
	request := &types.GetQuoteRequestData{
		ID:        quoteID,
		AccountID: accountID,
	}

	fqd := &types.FullQuoteData{
		ID: 1234,
		Req: &types.CreateQuoteRequestData{
			SourceCurrency: "BYN",
			TargetCurrency: "EUR",
			Amount:         "2021",
			AccountID:      42,
		},
		Res: &types.CreateQuoteResponseData{
			QuoteID:        "",
			TransactionFee: fmt.Sprint(2021 * 0.07),
			EDT:            120,
		},
	}

	tests := []struct {
		name    string
		req     *types.GetQuoteRequestData
		prepare func(qr *mock.MockQuoteRepo)
		want    *types.FullQuoteData
		wantErr bool
	}{
		{
			"1. Error empty request",
			nil,
			func(qr *mock.MockQuoteRepo) {},
			nil,
			true,
		},
		{
			"2. Error on read quote",
			request,
			func(qr *mock.MockQuoteRepo) {
				qr.EXPECT().Read(ctx, fmt.Sprint(request.ID), fmt.Sprint(request.AccountID)).Return(nil, errors.New("some error"))
			},
			nil,
			true,
		},
		{
			"3. Success",
			request,
			func(qr *mock.MockQuoteRepo) {
				qr.EXPECT().Read(ctx, fmt.Sprint(request.ID), fmt.Sprint(request.AccountID)).Return(fqd, nil)
			},
			fqd,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			xSvcMock := mock.NewMockExternalSvcClient(ctrl)
			quoteRepoMock := mock.NewMockQuoteRepo(ctrl)

			tt.prepare(quoteRepoMock)

			s := &QuoteService{
				xSvc:      xSvcMock,
				quoteRepo: quoteRepoMock,
			}
			got, err := s.Read(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
