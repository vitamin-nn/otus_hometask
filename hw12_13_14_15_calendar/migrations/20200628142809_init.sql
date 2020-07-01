-- +goose Up
CREATE TABLE events (
    id serial primary key,
    title text,
    description text,
    start_at timestamptz not null,
    end_at timestamptz not null,
    notify_at timestamptz,
    user_id integer
);

CREATE INDEX ix_events_notify_at ON events USING btree (notify_at);
CREATE INDEX ix_events_user_start_at_end_at ON events (user_id, start_at, end_at);

-- +goose Down
DROP TABLE events;
DROP INDEX ix_events_notify_at;
DROP INDEX ix_events_user_start_at_end_at;
