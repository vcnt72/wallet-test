-- +goose Up
-- +goose StatementBegin
ALTER TABLE ledgers ADD COLUMN error_code VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
