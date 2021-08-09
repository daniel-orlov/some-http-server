package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"some-http-server/internal/types"
	"testing"
)
// TODO generate mock
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

	tests := []struct {
		name    string
		req     *types.CreateQuoteRequestData
		prepare func(xSvcMock *mock.xsm, quoteRepoMock *mock.qrm)
		want    *types.CreateQuoteResponseData
		wantErr bool
	}{
		// TODO fix test cases
		{
			"1. Error empty request",
			nil,
			func(xs *mock.xsm, qr *mock.qrm) {},
			nil,
			true,
		},
		{
			"2. Error on create quote",
			request,
			func(xs *mock.xsm, qr *mock.qrm) {
				xs.some()
			},
			nil,
			true,
		},
		{
			"3. Error on save quote",
			request,
			func(xs *mock.xsm, qr *mock.qrm) {

			},
			nil,
			true,
		},
		{
			"4. Success",
			request,
			func(xs *mock.xsm, qr *mock.qrm) {

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

			xSvcMock := mock.f()
			quoteRepoMock := mock.f2()

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
	request := &types.GetQuoteRequestData{
		ID: quoteID,
	}

	tests := []struct {
		name    string
		req     *types.GetQuoteRequestData
		prepare func()
		want    *types.FullQuoteData
		wantErr bool
	}{
		// TODO fix test cases
		{
			"1. Error empty request",
			nil,
			func(xs *mock.xsm, qr *mock.qrm) {},
			nil,
			true,
		},
		{
			"2. Error on read quote",
			request,
			func(xs *mock.xsm, qr *mock.qrm) {
				xs.some()
			},
			nil,
			true,
		},
		{
			"3. Success",
			request,
			func(xs *mock.xsm, qr *mock.qrm) {

			},
			&types.FullQuoteData{
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
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			xSvcMock := mock.f()
			quoteRepoMock := mock.f2()

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
