package repository

import (
	"context"
	"reflect"
	"some-http-server/internal/types"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func TestQuoteRepository_Save(t *testing.T) {
	ctx := context.Background()
	q := &types.FullQuoteData{
		"some_uuid",
		types.CreateQuoteRequestData{
			SourceCurrency: "BYN",
			TargetCurrency: "EUR",
			Amount:         "2021",
			AccountID:      42,
		},
		types.CreateQuoteResponseData{
			QuoteID:        "some_quote_id",
			TransactionFee: "141.47",
			EDT:            60,
		},
	}

	expectedQuery := `
	INSERT INTO quotes (.+)
	VALUES (.+)
	RETURNING id;`

	tests := []struct {
		name    string
		prepare func(db sqlmock.Sqlmock)
		want    string
		wantErr bool
	}{
		{
			"1. Error on save full quote data",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(q.AccountID, q.QuoteID, q.Amount, q.SourceCurrency, q.TargetCurrency, q.TransactionFee, q.EDT).
					WillReturnError(errors.New("some error"))
			},
			"", true,
		},
		{
			"2. Success on save full quote data",
			func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(q.ID)
				mock.ExpectQuery(expectedQuery).
					WithArgs(q.AccountID, q.QuoteID, q.Amount, q.SourceCurrency, q.TargetCurrency, q.TransactionFee, q.EDT).
					WillReturnRows(rows)
			},
			q.ID, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			defer func() {
				if expErr := mock.ExpectationsWereMet(); expErr != nil {
					t.Errorf("QuoteRepository.Save() there were unfulfilled expectations: %s", expErr)
				}
			}()

			tt.prepare(mock)

			r := NewQuoteRepo(sqlx.NewDb(db, "postgres"))
			got, err := r.Save(ctx, q)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuoteRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("QuoteRepository.Save() = mismatch {+want;-got}\n%s", diff)
			}
		})
	}
}

func TestQuoteRepository_Read(t *testing.T) {
	ctx := context.Background()
	accountID := "42"
	id := "some_uuid"


	q := &types.FullQuoteData{
		"some_uuid",
		types.CreateQuoteRequestData{
			SourceCurrency: "BYN",
			TargetCurrency: "EUR",
			Amount:         "2021",
			AccountID:      42,
		},
		types.CreateQuoteResponseData{
			QuoteID:        "some_quote_id",
			TransactionFee: "141.47",
			EDT:            60,
		},
	}

	expectedQuery := `
	SELECT id, account_id, quote_id, amount, source_currency, target_currency, transaction_fee, edt
	FROM quotes 
	WHERE id = (.+) AND account_id = (.+);`

	tests := []struct {
		name    string
		prepare func(db sqlmock.Sqlmock)
		want    *types.FullQuoteData
		wantErr bool
	}{
		{
			"1. Error on read full quote data",
			func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(expectedQuery).
					WithArgs(id, accountID).
					WillReturnError(errors.New("some error"))
			},
			nil, true,
		},
		{
			"2. Success on read full quote data",
			func(mock sqlmock.Sqlmock) {
				columns := []string{
					"id", "account_id", "quote_id", "amount", "source_currency", "target_currency", "transaction_fee", "edt",
				}
				rows := sqlmock.NewRows(columns).AddRow(q.ID, q.AccountID, q.QuoteID, q.Amount, q.SourceCurrency, q.TargetCurrency, q.TransactionFee, q.EDT)
				mock.ExpectQuery(expectedQuery).
					WithArgs(id, accountID).
					WillReturnRows(rows)
			},
			q, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			defer func() {
				if expErr := mock.ExpectationsWereMet(); expErr != nil {
					t.Errorf("QuoteRepository.Read() there were unfulfilled expectations: %s", expErr)
				}
			}()

			tt.prepare(mock)

			r := NewQuoteRepo(sqlx.NewDb(db, "postgres"))

			got, err := r.Read(ctx, id, accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuoteRepository.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuoteRepository.Read() got = %v, want %v", got, tt.want)
			}
		})
	}
}
