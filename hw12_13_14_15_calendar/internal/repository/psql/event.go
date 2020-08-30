package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgconn"

	// pg driver.
	_ "github.com/jackc/pgx/v4/stdlib"
	outErr "github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/error"
	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

const (
	Type                    = "psql"
	ConstraintViolationCode = "23"
)

var _ repository.EventRepo = (*Psql)(nil)

type Psql struct {
	dsn string
	db  *sql.DB
}

func NewEventRepo(dsn string) *Psql {
	return &Psql{
		dsn: dsn,
	}
}

func (e *Psql) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", e.dsn)
	if err != nil {
		return err
	}
	e.db = db
	e.db.Stats()

	return e.db.PingContext(ctx)
}

func (e *Psql) Close() error {
	return e.db.Close()
}

func (e *Psql) CreateEvent(ctx context.Context, event *repository.Event) (*repository.Event, error) {
	var eventID int
	var notifyAt *time.Time
	if !event.NotifyAt.IsZero() {
		notifyAt = &event.NotifyAt
	}

	err := e.db.QueryRowContext(
		ctx,
		`INSERT INTO events(
			title,
			description,
			during,
			notify_at,
			user_id
		)
		VALUES($1, $2, tstzrange($3, $4, '[]'), $5, $6)
		RETURNING id`,
		event.Title,
		event.Description,
		event.StartAt,
		event.EndAt,
		notifyAt,
		event.UserID,
	).Scan(&eventID)
	if err != nil {
		specErr := getSpecificError(err)
		if specErr == nil {
			specErr = fmt.Errorf("insert error: %v", err)
		}

		return nil, specErr
	}

	event.ID = eventID

	return event, nil
}

func (e *Psql) UpdateEvent(ctx context.Context, event *repository.Event) (*repository.Event, error) {
	var notifyAt *time.Time
	if !event.NotifyAt.IsZero() {
		notifyAt = &event.NotifyAt
	}

	_, err := e.db.ExecContext(
		ctx,
		`UPDATE events
		 SET
			title = $1,
			description = $2,
			during = tstzrange($3, $4, '[]'),
			notify_at = $5,
			user_id = $6,
			notification_sent = $7
		 WHERE id = $8`,
		event.Title,
		event.Description,
		event.StartAt,
		event.EndAt,
		notifyAt,
		event.UserID,
		event.NotificationSent,
		event.ID,
	)
	if err != nil {
		specErr := getSpecificError(err)
		if specErr == nil {
			specErr = fmt.Errorf("update error: %v", err)
		}

		return nil, specErr
	}

	return event, nil
}

func (e *Psql) DeleteEvent(ctx context.Context, eventID int) error {
	res, err := e.db.ExecContext(
		ctx,
		"DELETE FROM events WHERE id = $1",
		eventID,
	)
	if err != nil {
		return fmt.Errorf("delete from events error: %v", err)
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt < 1 {
		return outErr.ErrEventNotFound
	}

	return nil
}

func (e *Psql) GetEventByID(ctx context.Context, eventID int) (*repository.Event, error) {
	var notifyAt sql.NullTime
	var notificationSent sql.NullTime

	ev := new(repository.Event)
	err := e.db.QueryRowContext(
		ctx,
		`SELECT
			id,
			title,
			description,
			lower(during),
			upper(during),
			notify_at,
			user_id,
			notification_sent
		FROM events
		WHERE
			id = $1
		LIMIT 1`,
		eventID,
	).Scan(&ev.ID,
		&ev.Title,
		&ev.Description,
		&ev.StartAt,
		&ev.EndAt,
		&notifyAt,
		&ev.UserID,
		&notificationSent,
	)
	if err != nil {
		return nil, fmt.Errorf("error while get event: %v", err)
	}

	if notifyAt.Valid {
		ev.NotifyAt = notifyAt.Time
	}

	if notificationSent.Valid {
		ev.NotificationSent = notificationSent.Time
	}

	if err != nil {
		return nil, fmt.Errorf("select exists error: %v", err)
	}

	return ev, nil
}

func (e *Psql) GetEventsByFilter(ctx context.Context, userID int, begin time.Time, end time.Time) ([]*repository.Event, error) {
	rows, err := e.db.QueryContext(
		ctx,
		`SELECT
			id,
			title,
			description,
			lower(during),
			upper(during),
			notify_at,
			user_id,
			notification_sent
		FROM events
		WHERE
			user_id = $1
			AND during && tstzrange($2, $3, '[]')`,
		userID,
		begin,
		end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*repository.Event

	for rows.Next() {
		ev := new(repository.Event)

		var notifyAt sql.NullTime
		var notificationSent sql.NullTime

		if err := rows.Scan(
			&ev.ID,
			&ev.Title,
			&ev.Description,
			&ev.StartAt,
			&ev.EndAt,
			&notifyAt,
			&ev.UserID,
			&notificationSent,
		); err != nil {
			return nil, fmt.Errorf("error while scan: %v", err)
		}

		if notifyAt.Valid {
			ev.NotifyAt = notifyAt.Time
		}
		if notificationSent.Valid {
			ev.NotificationSent = notificationSent.Time
		}

		result = append(result, ev)
	}

	return result, rows.Err()
}

func (e *Psql) DeleteOldEvents(ctx context.Context, t time.Time) error {
	_, err := e.db.ExecContext(
		ctx,
		"DELETE FROM events WHERE upper(during) < $1",
		t,
	)
	if err != nil {
		return fmt.Errorf("delete old events error: %v", err)
	}

	return nil
}

func (e *Psql) GetNotifyEvents(ctx context.Context, t time.Time) ([]*repository.Notification, error) {
	rows, err := e.db.QueryContext(
		ctx,
		`SELECT
			id,
			title,
			lower(during),
			user_id
		FROM events
		WHERE
			notify_at < $1
			AND notification_sent is null`,
		t,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*repository.Notification

	for rows.Next() {
		n := new(repository.Notification)

		if err := rows.Scan(
			&n.EventID,
			&n.EventTitle,
			&n.StartAt,
			&n.NotifyUserID,
		); err != nil {
			return nil, fmt.Errorf("error while scan: %v", err)
		}

		result = append(result, n)
	}

	return result, rows.Err()
}

func (e *Psql) UpdateNotificationSent(ctx context.Context, eventID int, t time.Time) error {
	_, err := e.db.ExecContext(
		ctx,
		`UPDATE events
		 SET
			notification_sent = $1
		 WHERE id = $2`,
		t,
		eventID,
	)
	if err != nil {
		specErr := getSpecificError(err)
		if specErr == nil {
			specErr = fmt.Errorf("update error: %v", err)
		}

		return specErr
	}

	return nil
}

func getSpecificError(err error) error {
	if errPg, ok := err.(*pgconn.PgError); ok {
		if sqlState := errPg.SQLState(); len(sqlState) > 2 && sqlState[0:2] == ConstraintViolationCode {
			return outErr.ErrDateBusy
		}
	}

	return nil
}
