-- +goose Up
-- +goose StatementBegin
CREATE TABLE ledgers(
  id BIGSERIAL PRIMARY KEY,
  wallet_id bigint not null,
  "type" varchar not null,
  idempotency_key varchar not null unique,
  amount bigint not null CHECK (amount >= 0),
  status varchar not null,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  CONSTRAINT fk_wallets FOREIGN KEY (wallet_id) REFERENCES wallets(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ledgers;
-- +goose StatementEnd
