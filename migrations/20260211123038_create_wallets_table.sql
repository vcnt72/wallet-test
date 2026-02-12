-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets(
  id BIGSERIAL PRIMARY KEY,
  user_id bigint not null,
  balance bigint not null check (balance >= 0),
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp,
  CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
