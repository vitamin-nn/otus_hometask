package psql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/vitamin-nn/otus_hometask/hw12_13_14_15_calendar/internal/repository"
)

const Type = "psql"

var _ repository.EventRepo = (*Psql)(nil)

type Psql struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *Psql {
	return &Psql{
		db: db,
	}
}

func (e *Psql) CreateEvent(ctx context.Context, event *repository.Event) (*repository.Event, error) {
	tx, err := e.db.Begin()
	if err != nil {
		err = fmt.Errorf("begin transaction error: %v", err)
		return nil, err
	}
	defer func() {
		err = tx.Rollback()
	}()

	exists, err := e.isBusyTime(ctx, event.UserID, event.StartAt, event.EndAt)
	if err != nil {
		return nil, err
	}
	if exists {
		err = repository.ErrDateBusy
		return nil, err
	}

	result, err := tx.ExecContext(
		ctx,
		`insert into events(
			title,
			description,
			start_at,
			end_at,
			notify_at,
			user_id
		)
		values($1, $2, $3, $4, $5, $6)`,
		event.Title,
		event.Description,
		event.StartAt,
		event.EndAt,
		event.NotifyAt,
		event.UserID,
	)
	if err != nil {
		err = fmt.Errorf("insert into events error: %v", err)
		return nil, err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		err = fmt.Errorf("last insert id getting error: %v", err)
		return nil, err
	}
	event.ID = int(eventID)

	err = tx.Commit()
	return event, err
}

func (e *Psql) UpdateEvent(ctx context.Context, eventID int, event *repository.Event) (*repository.Event, error) {
	tx, err := e.db.Begin()
	if err != nil {
		err = fmt.Errorf("begin transaction error: %v", err)
		return nil, err
	}
	defer func() {
		err = tx.Rollback()
	}()

	events, err := e.getEventsByFilter(ctx, event.UserID, event.StartAt, event.EndAt)
	if err != nil {
		err = fmt.Errorf("get events error: %v", err)
		return nil, err
	}
	if len(events) > 1 || (len(events) == 1 && events[0].ID != eventID) {
		err = repository.ErrDateBusy
		return nil, err
	}

	_, err = tx.ExecContext(
		ctx,
		`update events 
		 set
			title = $1,
			description = $2,
			start_at = $3,
			end_at = $4,
			notify_at = $5,
			user_id = $6
		 where id = $7`,
		event.Title,
		event.Description,
		event.StartAt,
		event.EndAt,
		event.NotifyAt,
		event.UserID,
		eventID,
	)
	if err != nil {
		err = fmt.Errorf("update events error: %v", err)
		return nil, err
	}
	event.ID = eventID
	err = tx.Commit()
	return event, err
}

func (e *Psql) DeleteEvent(ctx context.Context, eventID int) error {
	res, err := e.db.ExecContext(
		ctx,
		"delete from events where id = $1",
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
		return repository.ErrEventNotFound
	}

	return nil
}
func (e *Psql) GetEventsDay(ctx context.Context, userID int, dBegin time.Time) ([]*repository.Event, error) {
	year, month, day := dBegin.Date()
	loc := dBegin.Location()
	begin := time.Date(year, month, day, 0, 0, 0, 0, loc)
	end := time.Date(year, month, day, 23, 59, 59, 0, loc)
	return e.getEventsByFilter(ctx, userID, begin, end)
}

func (e *Psql) GetEventsWeek(ctx context.Context, userID int, wBegin time.Time) ([]*repository.Event, error) {
	end := wBegin.AddDate(0, 0, 7)
	return e.getEventsByFilter(ctx, userID, wBegin, end)
}

func (e *Psql) GetEventsMonth(ctx context.Context, userID int, mBegin time.Time) ([]*repository.Event, error) {
	end := mBegin.AddDate(0, 1, 0)
	return e.getEventsByFilter(ctx, userID, mBegin, end)
}

func (e *Psql) GetEventByID(ctx context.Context, eventID int) (*repository.Event, error) {
	var notifyAt sql.NullTime
	ev := new(repository.Event)
	err := e.db.QueryRowContext(
		ctx,
		`select
			id,
			title,
			description,
			start_at,
			end_at,
			notify_at,
			user_id
		from events
		where
			id = $1
		limit 1`,
		eventID,
	).Scan(&ev.ID,
		&ev.Title,
		&ev.Description,
		&ev.StartAt,
		&ev.EndAt,
		&notifyAt,
		&ev.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("error while get event: %v", err)
	}

	if notifyAt.Valid {
		ev.NotifyAt = notifyAt.Time
	}

	if err != nil {
		return nil, fmt.Errorf("select exists error: %v", err)
	}

	return ev, nil
}

func (e *Psql) getEventsByFilter(ctx context.Context, userID int, begin time.Time, end time.Time) ([]*repository.Event, error) {
	rows, err := e.db.QueryContext(
		ctx,
		`select
			id,
			title,
			description,
			start_at,
			end_at,
			notify_at,
			user_id
		from events
		where
			user_id = $1
			and (
				(start_at <= $2 and end_at >= $2)
				or (start_at <= $3 and end_at >= $3)
				or (start_at >= $2 and end_at <= $3)
			)`,
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

		if err := rows.Scan(
			&ev.ID,
			&ev.Title,
			&ev.Description,
			&ev.StartAt,
			&ev.EndAt,
			&notifyAt,
			&ev.UserID,
		); err != nil {
			return nil, fmt.Errorf("error while scan: %v", err)
		}

		if notifyAt.Valid {
			ev.NotifyAt = notifyAt.Time
		}
		result = append(result, ev)
	}

	return result, rows.Err()
}

func (e *Psql) isBusyTime(ctx context.Context, userID int, begin time.Time, end time.Time) (bool, error) {
	var exists bool
	err := e.db.QueryRowContext(
		ctx,
		`select exists(
			select 1 from events 
			where user_id = $1
			and (
				(start_at <= $2 and end_at >= $2)
				or (start_at <= $3 and end_at >= $3)
				or (start_at >= $2 and end_at <= $3)
			) 
			limit 1
		)`,
		userID,
		begin,
		end,
	).Scan(&exists)
	if err != nil {
		return exists, fmt.Errorf("select exists error: %v", err)
	}

	return exists, nil
}
