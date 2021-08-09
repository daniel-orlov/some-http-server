CREATE TABLE IF NOT EXISTS quotes (
    id BIGINT GENERATED ALWAYS AS IDENTITY NOT NULL,
    account_id BIGINT NOT NULL,
    quote_id VARCHAR(64) NOT NULL,
    amount VARCHAR(64) NOT NULL,
    source_currency VARCHAR(3) NOT NULL,
    target_currency VARCHAR(3) NOT NULL,
    transaction_fee VARCHAR(64) NOT NULL,
    edt BIGINT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS quotes_acc_id_quote_id_uniq_idx
ON quotes (account_id, quote_id);
