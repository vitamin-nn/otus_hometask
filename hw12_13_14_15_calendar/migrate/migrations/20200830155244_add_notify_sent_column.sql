-- +goose Up
-- +goose StatementBegin
ALTER TABLE events ADD COLUMN notification_sent timestamptz;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events DROP COLUMN IF EXISTS notification_sent;
-- +goose StatementEnd
