-- +goose Up
-- +goose StatementBegin
ALTER TABLE private_message 
ADD COLUMN deleted BOOLEAN DEFAULT FALSE,
ADD COLUMN updated BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE private_message 
DROP COLUMN deleted,
DROP COLUMN updated;
-- +goose StatementEnd