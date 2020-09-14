-- +goose Up
-- перед запуском нужно убедиться, что установлено расширение CREATE EXTENSION btree_gist; - нужны админские права
CREATE TABLE events (
    id serial primary key,
    title text,
    description text,
    during tstzrange,
    notify_at timestamptz,
    user_id integer,
    EXCLUDE USING GIST (user_id WITH =, during WITH &&)
);

CREATE INDEX ix_events_notify_at ON events USING btree (notify_at);

-- +goose Down
DROP TABLE events;
DROP INDEX ix_events_notify_at;
DROP INDEX ix_events_user_start_at_end_at;
