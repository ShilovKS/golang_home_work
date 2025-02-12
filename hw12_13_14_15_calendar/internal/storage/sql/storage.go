package sqlstorage

import (
	"context"
	"errors"

	"github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // драйвер postgres
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("the specified date is already occupied by another event")
)

type Storage struct {
	db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context, dsn string, driver string) error {
	db, err := sqlx.ConnectContext(ctx, driver, dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	query := `
        SELECT COUNT(*) 
        FROM events
        WHERE start_time < $2 AND end_time > $1
    `
	var count int
	selectErr := s.db.GetContext(ctx, &count, query, event.EndTime, event.StartTime)
	if selectErr != nil {
		return selectErr
	}
	if count > 0 {
		return ErrDateBusy
	}

	addQuery := `
        INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO NOTHING
    `
	_, err := s.db.ExecContext(ctx,
		addQuery,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.NotifyBefore)
	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	query := `
        UPDATE events
        SET title = $1, start_time = $2, end_time = $3, description = $4, notify_before = $5
        WHERE id = $6
    `
	result, err := s.db.ExecContext(ctx,
		query,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.NotifyBefore,
		id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `
        DELETE FROM events WHERE id = $1
    `
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrEventNotFound
	}

	return nil
}

func (s *Storage) ListEvents(ctx context.Context) ([]storage.Event, error) {
	query := `
        SELECT id, title, start_time, end_time, description, user_id, notify_before, created_at
        FROM events
    `

	var events []storage.Event
	err := s.db.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}
