-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
  id BIGSERIAL PRIMARY KEY,
  name varchar not null,
  created_at timestamptz default current_timestamp,
  updated_at timestamptz default current_timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
