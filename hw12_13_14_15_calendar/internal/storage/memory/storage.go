package memorystorage

import (
	"context"
	"errors"
	"sync"

	"github.com/ShilovKS/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("the specified date is already occupied by another event")
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) AddEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existingEvent := range s.events {
		if event.StartTime.Before(existingEvent.EndTime) && existingEvent.StartTime.Before(event.EndTime) {
			return ErrDateBusy
		}
	}

	if event.ID == "" {
		event.ID = generateUUID()
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.events[id]
	if !exists {
		return ErrEventNotFound
	}

	for _, e := range s.events {
		if e.ID == id {
			continue
		}
		if (e.StartTime.Before(event.EndTime) || e.StartTime.Equal(event.EndTime)) &&
			(e.EndTime.After(event.StartTime) || e.EndTime.Equal(event.StartTime)) {
			return ErrDateBusy
		}
	}

	event.ID = existing.ID
	s.events[id] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) ListEvents(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]storage.Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}

	return events, nil
}

func generateUUID() string {
	return uuid.New().String()
}
