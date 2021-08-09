package repository

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"some-http-server/internal/types"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type QuoteRepository struct {
	db *sqlx.DB
}

func NewQuoteRepo(db *sqlx.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

const saveQuoteQuery = `
	INSERT INTO quotes (account_id, quote_id, amount, source_currency, target_currency, transaction_fee, edt)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id
;`

func (r *QuoteRepository) Save(ctx context.Context, q *types.FullQuoteData) (string, error) {
	if q == nil {
		return "", errors.New("empty quote data")
	}

	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"quote": q,
	})

	var id string
	err := r.db.GetContext(ctx, &id, saveQuoteQuery, q.Req.AccountID, q.Res.QuoteID, q.Req.Amount, q.Req.SourceCurrency, q.Req.TargetCurrency, q.Res.TransactionFee, q.Res.EDT)
	if errors.Is(err, sql.ErrNoRows) {
		log.Debug("quote already exists for this user")
		return "", nil
	}
	if err != nil {
		return "", errors.Wrap(err, "cannot save quote in db")
	}

	return id, nil
}

const readQuoteQuery = `
	SELECT id, account_id, quote_id, amount, source_currency, target_currency, transaction_fee, edt
	FROM quotes 
	WHERE id = $1 AND account_id = $2;
	`

func (r *QuoteRepository) Read(ctx context.Context, id, accountID string) (*types.FullQuoteData, error) {
	if id == "" {
		return nil, errors.New("empty quote id")
	}

	if accountID == "" {
		return nil, errors.New("empty account id")
	}

	log := ctxlogrus.Extract(ctx).WithFields(logrus.Fields{
		"quote_id":   id,
		"account_id": accountID,
	})

	var fqd types.FullQuoteData
	err := r.db.GetContext(ctx, &fqd, readQuoteQuery, id, accountID)
	if errors.Is(err, sql.ErrNoRows) {
		log.Debug("quote with this account_id and id pair doesn't exist")
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "cannot read quote from db")
	}

	return &fqd, nil
}
