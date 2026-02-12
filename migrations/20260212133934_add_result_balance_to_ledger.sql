-- +goose Up
-- +goose StatementBegin
ALTER TABLE ledgers ADD COLUMN result_balance BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
